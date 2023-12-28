// const requestURL = "https://jsonplaceholder.typicode.com/users";
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
    console.log(data);
    console.log(window.location.pathname);
    
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
  //Конец

  //GET REQUEST
  function sendGETRequest(url, param1, param2) {
    var filter1 = document.getElementById(param1).value;
    var limit = document.getElementById('Limit').value;
    q = 'http://localhost:8080/api/' + url + param1 +'='+ encodeURIComponent(filter1)
    if (param2 !== undefined) {
      var filter2 = document.getElementById(param2).value;
      q = q + '&' + param2 +'='+ filter2
    }
    if (limit !== '') {
      q = q+'&Limit='+limit
    }
    
    fetch(q)
        .then(response => response.json())
        .then(data => {
            console.log('Ответ от сервера:', data);
            if (data.list.length != 0) {
              console.log(data)
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
  
    var day = date.getDate().toString().padStart(2, '0');
    var month = (date.getMonth() + 1).toString().padStart(2, '0');
    var year = date.getFullYear();
  
    var hours = date.getHours().toString().padStart(2, '0');
    var minutes = date.getMinutes().toString().padStart(2, '0');
    var seconds = date.getSeconds().toString().padStart(2, '0');
  
    return `${day}.${month}.${year} ${hours}:${minutes}:${seconds}`;
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

  //Успешное добавление
  function success() {
  Swal.fire({
    title: "Good job!",
    text: "You clicked the button!",
    icon: "success"
  });
}