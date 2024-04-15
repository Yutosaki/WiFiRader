// todo
// ・getTimeSlotからいれた予定の呼び出し方
// ・日付またぎの予定の考慮
// -・calculateShiftsのdayのアライメント(現状は15分区切りですべてhourで扱っている)
// ・preferredTimeの考慮(最悪やらなくてもいいと思う)

document.getElementById('submit').addEventListener('click',async () => {
	handleAuthClick(); // callbackでgen_scheduleを呼び出している
});

async function gen_schedule(){

	// formから数値を取得
	let minWeeklyHours = Number(document.getElementById('min-weekly-hours').value);
	let minShiftHours = Number(document.getElementById('min-hours').value);
	let transferTime = Number(document.getElementById('transfer-time').value) / 60;

	let startTime_string = document.getElementById('basaki-stime').value;
	const [hours, minutes] = startTime_string.split(':');
	let startTime = parseInt(hours, 10) + parseInt(minutes, 10) / 60;

	let endTime_string = document.getElementById('basaki-etime').value;
	const [hours2, minutes2] = endTime_string.split(':');
	let endTime = parseInt(hours2, 10) + parseInt(minutes2, 10) / 60;

	// eventsのTrasferTimeを初期化
	initEventsTrasferTime(transferTime);

	console.log(events);

	// google calenderから予定を呼び出す
	let existingEvents = await read_google_calender();
	// console.log(existingEvents);

	// getTimeSlotからいれた予定を呼び出す
	// 現状呼び出し方が不明のためスキップ

	// calculateShiftsを呼び出してシフトを作成する
	let shifts = calculateShifts(minWeeklyHours, minShiftHours, transferTime, startTime, endTime, existingEvents);
	if (shifts === null){
		alert("エラーが発生しました。もう一度決定ボタンを押してください。繰り返し表示される場合、週の最低労働時間が大きすぎたり、一出勤の最低労働時間が小さすぎたりしないか、ご確認ください。");
		return;
	}
	console.log(events);
}

function initEventsTrasferTime(transferTime){
	for(let i = 1; i <= 7; i++){
		let eventObject = events.find(event => event.dayOfWeek === i);
		eventObject.commuteTime = convertDecimalToTime(transferTime);
	}
}

async function read_google_calender(){
	let response;
	let now = new Date();
	let dayOfWeek = now.getDay(); // 0 (Sunday) to 6 (Saturday)
	let daysUntilNextMonday = (dayOfWeek === 0) ? 1 : (8 - dayOfWeek);
	let nextWeekStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() + daysUntilNextMonday);
	let nextWeekEnd = new Date(now.getFullYear(), now.getMonth(), now.getDate() + daysUntilNextMonday + 6);
	const request = {
		'calendarId': 'primary',
		'timeMin': nextWeekStart.toISOString(),
		'timeMax': nextWeekEnd.toISOString(),
		'showDeleted': false,
		'singleEvents': true,
		'maxResults': 10,
		'orderBy': 'startTime',
	};
	response = await gapi.client.calendar.events.list(request);
	cal_events = response.result.items;

	// 現状日付をまたぐ予定は考慮していない
	let schedule = cal_events
    .filter(event => event.start.dateTime && event.end.dateTime)
    .reduce((acc, event) => {
        let start = new Date(event.start.dateTime);
        let end = new Date(event.end.dateTime);
        let day = start.getDay(); // 0 (Sunday) to 6 (Saturday)
        if (day === 0){
            day = 7; // Sundayを7に変換
        }
        let time = {
            start: start.getHours() + start.getMinutes() / 60,
            end: end.getHours() + end.getMinutes() / 60
        };
		// eventsにカレンダーの情報を追加
		pushOboject(day, convertDecimalToTime(time.start), convertDecimalToTime(time.end), "googleCalendarEvent");

        // 同じdayのオブジェクトがあるかチェック
        let dayObject = acc.find(obj => obj.day === day);
        if (dayObject) {
            // 既に同じdayのオブジェクトがある場合、timeを追加
            dayObject.time.push(time);
        } else {
            // 同じdayのオブジェクトがない場合、新しいオブジェクトを作成
            acc.push({
                day: day,
                time: [time]
            });
        }
        return acc;
    }, []);
	return schedule;
}


// 少なくともこの関数内ではtransfer はhour
// アライメントはしてない
function calculateShifts(minWeeklyHours, minShiftHours, transferTime, startTime, endTime, existingEvents) {
    // 1週間分のシフトを格納する配列
	let possibleShifts = [];

    // 週の各日についてシフトを計算
	// day = 0 (Monday) to 6 (Sunday)
    for(let day = 0; day < 7; day++) {
        // その日の予定を取得
        let dailyEvents = existingEvents.filter(event => event.day - 1 === day);

        // その日の可能なシフトを計算
        for(let i = startTime; i <= endTime - minShiftHours; i+=0.25) {
			for(let j = i + minShiftHours; j <= endTime; j+=0.25){
				possibleShifts.push({day: day + 1, start: i, end: j});
			}
        }

        // 既存の予定と重複するシフトを除外
        dailyEvents.forEach(event => {
            let eventStart = event.start;
            let eventEnd = event.end;
            possibleShifts = possibleShifts.filter(possibleShift => 
				(possibleShift.end <= eventStart - transferTime || possibleShift.start >= eventEnd + transferTime));
        });
	}

	

    // 週最低労働時間を満たすまでランダムなシフトを追加
    let totalHours = 0;
    while(totalHours < minWeeklyHours) {
        // ランダムなインデックスを生成
        let randomIndex = Math.floor(Math.random() * possibleShifts.length);
		if (possibleShifts.length === 0){
			return null;
		}
        let shift = possibleShifts[randomIndex];

        // シフトを追加
		pushOboject(shift.day, convertDecimalToTime(shift.start), convertDecimalToTime(shift.end), "shift");
		totalHours += shift.end - shift.start;

        // 同じシフトを選ばないように、選択したシフトを可能なシフトから削除
        possibleShifts.splice(randomIndex, 1);

		// かぶる時間のシフトを削除
		possibleShifts = possibleShifts.filter(possibleShift => 
			possibleShift.day !== shift.day || (possibleShift.day === shift.day
			&& (possibleShift.end <= shift.start || possibleShift.start >= shift.end)));
    }
    
    return events;
}

function convertDecimalToTime(decimalTime) {
    let hours = Math.floor(decimalTime);
    let minutes = Math.round((decimalTime - hours) * 60);
    return hours.toString().padStart(2, '0') + ':' + minutes.toString().padStart(2, '0');
}

function pushOboject(day, start, end, type){
	if (type === "shift"){
		let shiftObject = events.find(event => event.dayOfWeek === day);
        shiftObject.shifts.push({start: start, end: end});
	}else if (type === "googleCalendarEvent"){
		let googleCalendarEventObject = events.find(event => event.dayOfWeek === day);
		googleCalendarEventObject.googleCalendarEvents.push({start: start, end: end});
	}else if (type === "class"){
		let classObject = events.find(event => event.dayOfWeek === day);
		classObject.classes.push({start: start, end: end});
	}else{
		alert("typeが不正です。書いているコードを確認してください。");
	}
}