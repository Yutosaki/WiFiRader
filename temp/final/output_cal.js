// var events = [
//   {
//     dayOfWeek: 1, // 1は月曜日
//     commuteTime: "00:15", // 通勤にかかる時間
//     classes: [
//       { start: "05:00", end: "08:10" },
//       { start: "10:20", end: "16:50" }
//     ],
//     shifts: [
//       { start: "03:00", end: "4:10" },
//       { start: "18:20", end: "23:50" }
//     ],
//     googleCalendarEvents: [ // Google カレンダーのイベントを配列で管理
//       { start: "16:30", end: "18:00" },
//       { start: "1:00", end: "2:00" } 
//     ]
//   },
//   {
//     dayOfWeek: 2, // 2は火曜日
//     commuteTime: "00:15", // 通勤にかかる時間
//     classes: [
//       { start: "05:00", end: "08:10" },
//       { start: "18:20", end: "23:50" }
//     ],
//     shifts: [
//       { start: "05:00", end: "08:10" },
//       { start: "18:20", end: "23:50" }
//     ],
//     googleCalendarEvents: [ // Google カレンダーのイベントを配列で管理
//       { start: "16:30", end: "18:00" },
//       { start: "1:00", end: "3:00" } 
//     ]
//   },
//   {
//     dayOfWeek: 3, // 2は火曜日
//     commuteTime: "00:15", // 通勤にかかる時間
//     shifts: [
//       { start: "05:00", end: "08:10" },
//       { start: "18:20", end: "23:50" }
//     ]
//   }
// ];

document.addEventListener('DOMContentLoaded', () => {
    const outputContainer = document.querySelector('.output');
    initializeCalendar(outputContainer);
    generateTimeColumn(outputContainer);
    document.getElementById('generateSchedule').addEventListener('click', () => {
        fetchSchedule(outputContainer);
    });
  });
  
  function generateTimeColumn() {
      const timeColumn = document.getElementById('timeColumn');
      timeColumn.innerHTML = '';
      for (let hour = 0; hour < 24; hour++) {
          for (let minute = 0; minute < 60; minute +=5) {
              const timeDiv = document.createElement('div');
              timeDiv.classList.add('hour');
              if (minute === 0) {
                  timeDiv.textContent = `${String(hour).padStart(2, '0')}:00`;
              }
              timeColumn.appendChild(timeDiv);
          }
      }
  }
  
  function initializeCalendar() {
    const calendarContainer = document.getElementById('scheduleCalendar');
    calendarContainer.innerHTML = '';
    const days = ['月曜日', '火曜日', '水曜日', '木曜日', '金曜日', '土曜日', '日曜日'];
  
    days.forEach((day) => {
      const dayDiv = document.createElement('div');
      dayDiv.classList.add('day');
      const dayTitle = document.createElement('div');
      dayTitle.textContent = day;
      dayTitle.classList.add('dayTitle');
      dayDiv.appendChild(dayTitle);
  
      for (let hour = 0; hour < 24; hour++) {
        for (let minute = 0; minute < 60; minute += 5) {
          const minuteDiv = document.createElement('div');
          minuteDiv.classList.add('minute');
          if (minute === 30|| minute === 0) {
            minuteDiv.classList.add('half-hour');
          }
          dayDiv.appendChild(minuteDiv);
        }
      }
      calendarContainer.appendChild(dayDiv);
    });
  }
  
  async function fetchSchedule(outputContainer) {
    applyScheduleColors(events, outputContainer);
  }
  
  function timeToMinutes(time) {
    const [hours, minutes] = time.split(':').map(Number);
    return hours * 60 + minutes;
  }
  
  function applyScheduleColors(schedules, outputContainer) {
    const daysDiv = outputContainer.querySelectorAll('.day');
    daysDiv.forEach((dayDiv, index) => {
      const dayOfWeek = index + 1;
      schedules.forEach(schedule => {
        if (schedule.dayOfWeek === dayOfWeek) {
          const minuteDivs = dayDiv.querySelectorAll('.minute');
          console.log(schedule.classes);
          if (schedule.classes){
            schedule.classes.forEach(classEvent => {
              const classStartTime = timeToMinutes(classEvent.start);
              const classEndTime = timeToMinutes(classEvent.end);
              minuteDivs.forEach((minuteDiv, i) => {
                const currentTime = i * 5;
                if (currentTime >= classStartTime && currentTime < classEndTime) {
                  minuteDiv.classList.add('class-time');
                }
              });
            });
          }
          if (schedule.shifts){
            schedule.shifts.forEach(shift => {
              const shiftStartTime = timeToMinutes(shift.start);
              const shiftEndTime = timeToMinutes(shift.end);
              const commuteTimeMinutes = timeToMinutes(schedule.commuteTime);
              const commuteStartTime = shiftStartTime - commuteTimeMinutes;
              const commuteEndTime = shiftEndTime + commuteTimeMinutes;
    
              minuteDivs.forEach((minuteDiv, i) => {
                const currentTime = i * 5;
                // Apply work time
                if (currentTime >= shiftStartTime && currentTime < shiftEndTime) {
                  minuteDiv.classList.add('work-time');
                }
                // Apply commute time before and after the shift
                if ((currentTime >= commuteStartTime && currentTime < shiftStartTime) ||
                    (currentTime >= shiftEndTime && currentTime <= commuteEndTime)) {
                  minuteDiv.classList.add('commute-time');
                }
              });
            });
          }
          if (schedule.googleCalendarEvents){
            schedule.googleCalendarEvents.forEach(event => {
              const eventStartTime = timeToMinutes(event.start);
              const eventEndTime = timeToMinutes(event.end);
              minuteDivs.forEach((minuteDiv, i) => {
                const currentTime = i * 5;
                if (currentTime >= eventStartTime && currentTime < eventEndTime) {
                  minuteDiv.classList.add('google-calendar-time');
                }
              });
            });
          }
        }
      });
    });
  }
  