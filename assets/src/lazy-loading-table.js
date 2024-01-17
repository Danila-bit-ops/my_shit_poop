//Прослушивание скрола и при достижении 99 строки, вызов функции
document.addEventListener("DOMContentLoaded", function () {
  const tableContainer = document.querySelector(".table-container");
  var lastRow;
  var offset;
  var btn = document.getElementById("top");
  tableContainer.addEventListener("scroll", function () {
    var filter = window.fltr;
    var param1 = window.parametr1;
    var param2 = window.parametr2;
    const tbody = document.getElementById("table-tbody");
    const rows = tbody.getElementsByTagName("tr");
    const lastVisibleRow = getLastVisibleRow(tableContainer, rows);
    if (rows.length==100) {
      lastRow=99;
      offset=100;
    }

    if (lastVisibleRow>=lastRow){
      console.log(">"+lastRow);
      sendGETLazyRequest("lazy-loading",offset,filter,param1,param2);
      lastRow=lastRow+100;
      offset=offset+100;
    }
    if (lastVisibleRow >= 100) {
      btn.style.display = 'block';
    } else {
      btn.style.display = 'none';
    }
  });

  function getFirstVisibleRow(container, rows) {
    const containerTop = container.scrollTop;
    const rowHeight = rows[0].offsetHeight;
    return Math.floor(containerTop / rowHeight);
  }

  function getLastVisibleRow(container, rows) {
    const containerTop = container.scrollTop;
    const containerHeight = container.clientHeight;
    const rowHeight = rows[0].offsetHeight;

    const visibleRows = Math.ceil(containerHeight / rowHeight);
    const firstVisibleRow = getFirstVisibleRow(container, rows);

    return firstVisibleRow + visibleRows - 1;
  }
});

//Запрос на Golang для получения ещё 100 записей
function sendGETLazyRequest(url, offset, filter, param1, param2) {
  console.log(filter);
  if (filter==undefined) {
    q = 'http://localhost:8080/api/' + url + '?Offset='+ encodeURIComponent(offset)
  } else if (filter=="FilterParamID") {
    q = 'http://localhost:8080/api/' + url + '?' + filter + "=" + encodeURIComponent(param1) +'&Offset='+ encodeURIComponent(offset)
  } else {
    q = 'http://localhost:8080/api/' + url + '?RngStart=' + encodeURIComponent(param1) + '&RngEnd=' + encodeURIComponent(param2) + '&Offset=' + encodeURIComponent(offset)
  }
  // q = 'http://localhost:8080/api/' + url + '?Offset='+ encodeURIComponent(offset)
  fetch(q)
    .then(response => response.json())
    .then(data => {
      console.log('Ответ от сервера:', data);
      appendRowsToTableLazy(data)
      if (data.list.length != 0) {
        console.log(data);
        
      } else {
          // Вызов SweetAlert с уведомлением
          Swal.fire({
          icon: "error",
          title: "Упс...",
          text: "Больше значений нет!",
        });
      }
    })
    .catch(error => {
      console.error('Ошибка:', error);
    });
}
//Конец

//Вставить JSON в таблицу
function appendRowsToTableLazy(jsonData) {
  var tbody = document.getElementById("table-tbody");
  for (var i = 0; i < jsonData.list.length; i++) {
    var row = document.createElement("tr");
    var columnsOrder = ['ID', 'ParamID', 'Val', 'Timestamp', 'ChangeBy', 'XMLCreate', 'Manual', 'Comment'];
    for (var j = 0; j < columnsOrder.length; j++) {
      var key = columnsOrder[j];
      var cell = document.createElement("td");
      // Форматирование даты, если столбец 'Timestamp'
      if (key === 'Timestamp') {
        cell.textContent = formatDateTime(jsonData.list[i][key]);
      } else {
        cell.textContent = jsonData.list[i][key];
      }
      row.appendChild(cell);
    }
    tbody.appendChild(row);
  }
}
//Конец

//Формат даты
function formatDateTime(timestamp) {
  var date = new Date(timestamp);
  var formattedDateTime = date.toLocaleString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: 'UTC'
  });
  return formattedDateTime;
}