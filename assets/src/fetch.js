const requestURL = "http://localhost:8080/api/table-hour-params" ;
//Отправить запрос
function sendRequest(method, url, body = null) {
  const headers = {
    'Content-Type': 'application/json'
  }
  return fetch(url).then(response => {
    if (response.ok) {
      return response.json()
    } else {
      return response.json().then( error =>{
        const e = new Error('Something Wrong')
        e.data = error
        throw e
      })
    }
  })
}
//Конец
  
//Таблица при открытии index.html
if (window.location.pathname === '/index') {
  sendRequest('GET', requestURL)
  .then(data => {
    appendRowsToTable(data);
    animateTableContainer();
  })
  .catch(err => console.log(err));
}
//Конец

//Анимация таблицы
function animateTableContainer() {
  var tableContainer = document.querySelector(".table-container");
  tableContainer.classList.add("scroll-table-show"); 
}
//Конец FilterParamID

//GET REQUEST
function sendGETRequest(url, param1, param2) {
  scrollToTop();
  var filter1 = document.getElementById(param1).value;
  q = 'http://localhost:8080/api/' + url + param1 +'='+ encodeURIComponent(filter1)
  if (param2 !== undefined) {
    var filter2 = document.getElementById(param2).value;
    q = q + '&' + param2 +'='+ filter2
  }
  fetch(q)
    .then(response => response.json())
    .then(data => {
      console.log('Ответ от сервера:', data);
      if (data.list.length != 0) {
        window.fltr = param1;
        window.parametr1 = filter1;
        window.parametr2 = filter2;
        appendRowsToTable(data);

      } else {
          // Вызов SweetAlert с уведомлением
          Swal.fire({
          icon: "error",
          title: "Упс...",
          text: "Нет подходящего значения!",
        });
      }
    })
    .catch(error => {
      console.error('Ошибка:', error);
    }); 
}
//Конец

// POST запрос для удаления строки по ID
function deletePOSTRequest(param) {
  var filter = document.getElementById(param).value;
  var data = {
    ID: parseInt(filter)
  };
  q = 'http://localhost:8080/api/del-by-id';
  fetch(q,{
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data) // Отправляем JSON-объект
  })
  .then(response => response.json())
  .then(result => {
    if (result.error) {
      console.error('Ошибка:', result.error);
    } else if (result.message === "Данные успешно удалены из базы данных") {
      Swal.fire({
        title: "Успешно!",
        text: "Данные успешно удалены из базы данных!",
        icon: "success"
      });
    } else {
      Swal.fire({
        icon: "error",
        title: "Упс...",
        text: "Нет такого ID!",
      });
    }
  })
  .catch(error => console.error('Ошибка:', error));
}
//Конец

//POST запрос для добавления новой строки
function addNewRecord(value,paramid,timest,xml,manual,ch,comm) {
  var val = document.getElementById(value).value;
  var pid = document.getElementById(paramid).value;
  var timeString = document.getElementById(timest).value;
  var xmlcrText = document.getElementById(xml).value;
  var manText = document.getElementById(manual).value;
  var change = document.getElementById(ch).value;
  var comment = document.getElementById(comm).value;
    
  // Создаем объект data
  var data = {
    AddNewVal: val.toString(),
    AddNewParamID: pid.toString(),
    AddNewTimestamp: timeString.toString(),
    AddNewXml: xmlcrText.toString(),
    AddNewManual: manText.toString(),
    AddNewChange: change.toString(),
    AddNewComment: comment.toString()
  };
  q = 'http://localhost:8080/api/add-new-record';

  fetch(q,{
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data) // Отправляем JSON-объект
  })
  .then(response => response.json())
  .then(result => {
    if (result.error) {
      console.error('Ошибка:', result.error);
    } else {
      Swal.fire({
        title: "Успешно!",
        text: "Данные успешно добавлены в базу данных!",
        icon: "success"
      });
    }
  })
  .catch(error => console.error('Ошибка:', error));
}
//Конец

//POST запрос для изменения строки
function updRecord(id,value,paramid,timest,xml,manual,ch,comm) {
  var id = document.getElementById(id).value;
  var val = document.getElementById(value).value;
  var pid = document.getElementById(paramid).value;
  var timeString = document.getElementById(timest).value;
  var xmlcrText = document.getElementById(xml).value;
  var manText = document.getElementById(manual).value;
  var change = document.getElementById(ch).value;
  var comment = document.getElementById(comm).value;
    
  // Создаем объект data
  var data = {
    UpdID: id.toString(),
    UpdVal: val.toString(),
    UpdParamID: pid.toString(),
    UpdTimestamp: timeString.toString(),
    UpdXml: xmlcrText.toString(),
    UpdManual: manText.toString(),
    UpdChange: change.toString(),
    UpdComment: comment.toString()
  };

  q = 'http://localhost:8080/api/update-by-id';
  fetch(q,{
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data) // Отправляем JSON-объект
  })
  .then(response => response.json())
  .then(result => {
    if (result.error) {
      console.error('Ошибка:', result.error);
    } else if (result.message === "Данные успешно изменены") {
      Swal.fire({
        title: "Успешно!",
        text: "Данные успешно изменены!",
        icon: "success"
      });
    } else {
      Swal.fire({
        icon: "error",
        title: "Упс...",
        text: "Нет такого ID!",
      });
    }
  })
  .catch(error => console.error('Ошибка:', error));
}
//Конец

//Вставить JSON в таблицу
function appendRowsToTable(jsonData) {
  var tbody = document.getElementById("table-tbody");
  tbody.innerHTML = '';
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
//Конец
  
//Проверка ввода id & param_id
function validateInput(id,buttonId,errId) {
  var inputValue = document.getElementById(id).value;
  var inputStyle = document.getElementById(id).style
  var ErrImg = document.getElementById(errId).style 
  var submitButton = document.getElementById(buttonId);
  if (!isNaN(inputValue)) {
      submitButton.disabled = false;
      inputStyle.background = "#f1f1f1";
      ErrImg.display = "none"
  } else {
    submitButton.disabled = true;
    inputStyle.background = "#ffbd2c"
    ErrImg.display = "block"
  }
}
//Конец

//Проверка ввода val
function validateInputVal(id, buttonId) {
  var inputValue = document.getElementById(id).value;
  var inputStyle = document.getElementById(id).style
  var submitButton = document.getElementById(buttonId);
  if (!isNaN(inputValue) || (inputValue.startsWith('-') && !isNaN(inputValue.substring(1)))) {
    submitButton.disabled = false;
    inputStyle.background = "#f1f1f1";
  } else {
    submitButton.disabled = true;
    inputStyle.background = "#ffbd2c"
  }
}
//Конец

function scrollToTop() {
  const tableContainer = document.querySelector('.table-container');
  tableContainer.scrollTop = 0;
}

