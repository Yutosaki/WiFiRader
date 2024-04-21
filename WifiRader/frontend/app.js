document.addEventListener('DOMContentLoaded', function() {
    initMap();
});

function submitLocationAndAmount() {
    if (!navigator.geolocation) {
        console.error("Geolocation is not supported by your browser.");
        return;
    }

    const desiredAmount = document.getElementById('desiredAmount').value;
    if (!desiredAmount) {
        alert("金額を入力してください。");
        return;
    }

    navigator.geolocation.getCurrentPosition(position => {
        const pos = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
        };

        fetch('http://localhost:8080/submit-location', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ pos, desiredAmount })
        })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }, () => {
        alert("現在地の取得に失敗しました。");
    });
}

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

    initMap(); // Assuming this function initializes the Google Map
});

