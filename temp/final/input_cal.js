document.addEventListener('DOMContentLoaded', () => {
	initEvents();
	generateEmptyTimeSlots();
});

let startDragCell = '';
let startDragTime = '';
let startDragDay = '';
let endDragTime = '';
let endDragDay = '';
var events = [];

function initEvents(){
    for(let i = 1; i <= 7; i++){
        events.push({
            dayOfWeek: i,
            commuteTime: 0,
            classes: [],
            shifts: [],
            googleCalendarEvents: []
        });
    }
}

function generateTimeSlots() {
	var timeSlots = [];
	for (var hours = 0; hours < 24; hours++) {
		for (var minutes = 0; minutes < 60; minutes += 5) {
			var hourStr = (hours < 10 ? '0' : '') + hours;
			var minStr = (minutes < 10 ? '0' : '') + minutes;
			timeSlots.push(hourStr + ':' + minStr);
		}
	}
	return timeSlots;
}

function generateEmptyTimeSlots() {
	var calendarTable = document.querySelector('.calendar');
	var timeSlots = generateTimeSlots();
	var row = calendarTable.insertRow(0); // Insert new row at the end
	var i;
	var th = document.createElement('th');

	console.log("generateEmptyTimeSlots", events);
	th.classList.add('time-col');
	row.appendChild(th);
	for (i = 1; i < 8; i++){
		th = document.createElement('th');
		th.textContent = getDayName(i-1);
		row.appendChild(th);
	}
	// Loop through each time slot
	timeSlots.forEach(time => {
		row = calendarTable.insertRow(-1); // Insert new row at the end
		var timeCell = row.insertCell(0);
		timeCell.setAttribute("id", time);
		const regex = wildcardToRegex("*:00");
		if (regex.test(time)){
                    	var textContent = document.createElement('div');
			textContent.textContent = time;
			timeCell.appendChild(textContent);
			timeCell.classList.add('hour');
		// Loop through each day of the week
			for (var i = 1; i <= 7; i++) {
				var dayCell = row.insertCell(i);
				var dayName = getDayName(i - 1);
				dayCell.classList.add('empty-cell-border');
			}
		}else if(wildcardToRegex("*:30").test(time)){
			timeCell.classList.add('minute');
			// Loop through each day of the week
			for (var i = 1; i <= 7; i++) {
				var dayCell = row.insertCell(i);
				var dayName = getDayName(i - 1);
				dayCell.classList.add('empty-cell-border');
			}
		}else{
			timeCell.classList.add('minute');
			// Loop through each day of the week
			for (var i = 1; i <= 7; i++) {
				var dayCell = row.insertCell(i);
				var dayName = getDayName(i - 1);
				dayCell.classList.add('empty-cell');
			}
		}

	});
}

document.addEventListener('mousedown', (e) => {
	if (e.target.tagName === 'TD') {
		startDragCell = e.target;
		startDragTime = e.target.parentNode.cells[0].id;
		startDragDay = getDayName(e.target.cellIndex - 1);
	}
});

document.addEventListener('mouseup', (e) => {
	if (e.target.tagName === 'TD') {
		endDragTime = e.target.parentNode.cells[0].id;
		endDragDay = getDayName(e.target.cellIndex - 1);
		console.log('Dragged time range:', startDragDay, startDragTime, 'to',endDragDay, endDragTime);
	

		// Here you can add code to prompt the user for event details and add the event to the events array
		// For simplicity, we'll log the event datac
		if (startDragDay === endDragDay) {
			var newEvent = {
				day: startDragDay,
				sTime: startDragTime,
				eTime: endDragTime,
				title: prompt('Enter event title:')
			};
			// 全てのセルからハイライトを削除する
			removeAllCellHighlights();
			if (newEvent.title) {
				events.push(newEvent);
				pushOboject(getDayIndex(startDragDay), startDragTime, endDragTime, "class", newEvent.title);
				console.log('Event added:', newEvent);
				// Refresh the calendar to reflect the new event
					// document.querySelector('.calendar').innerHTML = ''; // Clear the calendar
				highlightEvents(); // Repopulate events
			}
		}
	}

	startDragTime = '';
	endDragTime = '';
});

document.addEventListener('mouseover', (e) => {
	if (e.target.tagName === 'TD' && startDragTime !== '') {
		const currentCellDay = getDayName(e.target.cellIndex - 1);
		if (currentCellDay === startDragDay) {
			currentTime = e.target.parentNode.cells[0].id;
			const startTime = startDragTime;
			
			const endRowIndex = e.target.parentNode.rowIndex;
			const startRowIndex = startDragCell.parentNode.rowIndex;
			
			//ドラッグされている先頭のセルに時間表示
			startDragCell.textContent = startTime + "~" + currentTime;
			
			// ドラッグされた時間範囲内のセルにクラスを追加する
			for (let i = startRowIndex; i <= endRowIndex; i++) {
				const row = document.querySelector('.calendar').rows[i];
				const timeValue = row.cells[0].id;
				if (timeValue >= startTime) {
					const eventCell = row.cells[e.target.cellIndex];
					eventCell.classList.add('cell-highlight');
				}
			}
		}
	}
});

// Function to highlight events in the calendar
function highlightEvents() {
	var calendarTable = document.querySelector('.calendar');
	
	// Loop through events array
	events.forEach(event => {
	    var startTime = event.sTime;
	    var endTime = event.eTime;
	    var dayIndex = getDayIndex(event.day);
	
	    // Loop through each time slot row in the calendar
	    for (var i = 1; i < calendarTable.rows.length; i++) {
	        var timeCell = calendarTable.rows[i].cells[0];
	        var timeValue = timeCell.id;
	
	        // Check if the time slot matches the event time range
	        if (timeValue >= startTime && timeValue <= endTime) {
	            var eventCell = calendarTable.rows[i].cells[dayIndex];
	            var textContent = document.createElement('div');
	            eventCell.classList.remove('empty-cell');
	            eventCell.classList.add('event-highlight'); // Add a CSS class for highlighting
	
	            if (timeValue == startTime) {
	                eventCell.classList.remove('event-highlight'); // Add a CSS class for highlighting
	           	eventCell.classList.add('event-name'); // Add a CSS class for highlighting
	                var closeBtn = document.createElement('button');
	                closeBtn.classList.add('close-btn');
	                closeBtn.innerHTML = 'x';
			closeBtn.onclick = () => {
				events = events.filter(e => e.day !== event.day || e.sTime !== event.sTime);
				calendarTable.innerHTML = '';
				generateEmptyTimeSlots();
				highlightEvents();
			};
			textContent.textContent = event.title;
	                calendarTable.rows[i+1].cells[dayIndex].appendChild(textContent);
	                eventCell.innerHTML = ''; // Clear the cell
	                eventCell.appendChild(closeBtn);
	            }
	        }
	    }
	});
}

function removeAllCellHighlights() {
	startDragCell.textContent = '';
	const highlightedCells = document.querySelectorAll('.cell-highlight');
	highlightedCells.forEach(cell => {
		cell.classList.remove('cell-highlight');
	});
}

function removeAllEventHighlights() {
	startDragCell.textContent = '';
	const highlightedCells = document.querySelectorAll('.event-highlight');
	highlightedCells.forEach(cell => {
		cell.classList.remove('event-highlight');
	});
}

function getDayName(dayIndex) {
	const daysOfWeek = [
		'月曜日',
		'火曜日',
		'水曜日',
		'木曜日',
		'金曜日',
		'土曜日',
		'日曜日'
	];
	return daysOfWeek[dayIndex];
}

function getDayIndex(dayName) {
	const daysOfWeek = [
		'月曜日',
		'火曜日',
		'水曜日',
		'木曜日',
		'金曜日',
		'土曜日',
		'日曜日'
	];
	return daysOfWeek.indexOf(dayName) + 1; // Adding 1 because of the time column in the calendar
}

function getTimeSlot(timeIndex) {
	const hours = Math.floor(timeIndex * 15 / 60); // 15分ごとの時間枠なので、60で割って時間に変換
	const minutes = (timeIndex * 15) % 60; // 15分ごとの時間枠の余りが分となる
	return
		`${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
}

function wildcardToRegex(wildcard) {
	// アスタリスク(*)を正規表現の「任意の文字列」に変換
	// クエスチョンマーク(?)を正規表現の「任意の1文字」に変換
	const regex = wildcard.replace(/\*/g, ".*").replace(/\?/g, ".");
	return new RegExp(`^${regex}$`);
}
