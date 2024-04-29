document.addEventListener('DOMContentLoaded', function() {
    initMap();
});

function initMap() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function(position) {
            var pos = {
                lat: position.coords.latitude,
                lng: position.coords.longitude
            };

            var map = new google.maps.Map(document.getElementById('map'), {
                center: pos,
                zoom: 15
            });

            new google.maps.Marker({
                position: pos,
                map: map,
                title: '現在地'
            });
        }, function() {
            handleLocationError(true, map, map.getCenter());
        });
    } else {
        handleLocationError(false, null, { lat: -34.397, lng: 150.644 });
    }
}

function submitLocationAndAmount() {
    if (!navigator.geolocation) {
        console.error("Geolocation is not supported by your browser.");
        return;
    }

    const amountInput = document.getElementById('desiredAmount').value;
    const desiredAmount = parseInt(amountInput, 10); // 文字列を整数に変換
    if (isNaN(desiredAmount)) {  // 数値変換が正しく行われたかチェック
        alert("金額を数値で入力してください。");
        return;
    }

    navigator.geolocation.getCurrentPosition(position => {
        const pos = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude
        };
        console.log("Position data:", pos);

        fetch('http://localhost:8080/submit-location', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ pos: pos, desiredAmount: desiredAmount})
        })
        .then(response => response.json())
        .then(data => {
            if (data !== null) { // nullチェックを追加
                console.log('Number of places received:', data.length);
                console.log("Position data:", pos);
                if (data.length > 0) {
                    addMarkerAndUrl(data, pos);
                } else {
                    console.error('No valid locations received:', data);
                    alert('カフェの情報が見つかりませんでした。');
                }
            } else {
                console.error('No valid locations received:', data);
                alert('カフェの情報が見つかりませんでした。');
            }
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }, () => {
        alert("現在地の取得に失敗しました。");
    });
}

function addMarkerAndUrl(places, pos) {
    console.log('Places Data:', places);
    const map = new google.maps.Map(document.getElementById('map'), {
        center: { lat: pos.latitude, lng: pos.longitude },
        zoom: 15
    });

    // 現在地のマーカー
    new google.maps.Marker({
        position: { lat: pos.latitude, lng: pos.longitude },
        map: map,
        icon: {
            path: google.maps.SymbolPath.CIRCLE,
            scale: 8,
            fillColor: 'blue',
            fillOpacity: 0.6,
            strokeColor: 'white',
            strokeWeight: 2
        },
        title: '現在地'
    });

    let nearestPlace = null;
    let shortestDistance = Infinity;

    const placePos = { lat: place.Latitude, lng: place.Longitude };
    places.forEach(place => {
        const distance = google.maps.geometry.spherical.computeDistanceBetween(
            { lat: pos.latitude, lng: pos.longitude },
            placePos
        );

        if (distance < shortestDistance) {
            shortestDistance = distance;
            nearestPlace = place;
        }
    });

    places.forEach(place => {
        const marker = new google.maps.Marker({
            position: placePos,
            map: map,
            icon: nearestPlace === place ? '' : 'http://maps.google.com/mapfiles/ms/icons/red-dot.png',
            title: place.name
        });

        const infowindow = new google.maps.InfoWindow({
            content: `<a href="${place.url}" target="_blank">${place.name}</a>`
        });

        marker.addListener('click', function() {
            infowindow.open(map, marker);
        });

        // サイドバーに最も近い場所の情報を表示
        if (nearestPlace === place) {
            console.log('Nearest place:', place);
            const placeIframe = document.getElementById('placeIframe');
            placeIframe.src = place.url ? place.url : "URLを取得できませんでした"; 
        }
    });
}

function handleLocationError(browserHasGeolocation, map, pos) {
    console.log(browserHasGeolocation ?
        'Error: The Geolocation service failed.' :
        'Error: Your browser does not support geolocation.');
    map = new google.maps.Map(document.getElementById('map'), {
        center: pos,
        zoom: 15
    });
    new google.maps.Marker({
        position: pos,
        map: map,
        title: 'Default Location'
    });
}

document.addEventListener('DOMContentLoaded', function() {
    const toggleButton = document.getElementById('toggleButton');
    const sidebar = document.getElementById('sidebar');
    let isSidebarOpen = true;

    toggleButton.addEventListener('click', function() {
        if (isSidebarOpen) {
            sidebar.style.transform = 'translateX(-100%)'; // サイドバーを左に隠す
        } else {
            sidebar.style.transform = 'translateX(0)'; // サイドバーを表示
        }
        isSidebarOpen = !isSidebarOpen;
    });

    initMap();//現在地の読み込み
});
