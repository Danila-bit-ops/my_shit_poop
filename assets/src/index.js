import App from "./App"

//Форма для добавления новой записи
function openForm() {
  document.getElementById("AddNewStr").style.display = "block";
  document.getElementById("UpdStr").style.display = "none";
  document.getElementById("ParamIDSearch").style.display = "none";
  document.getElementById("IDSearch").style.display = "none";
  document.getElementById("DeleteID").style.display = "none";
}
function closeForm() {
  document.getElementById("AddNewStr").style.display = "none";
}

function openForm1() {
  document.getElementById("UpdStr").style.display = "block";
  document.getElementById("AddNewStr").style.display = "none";
  document.getElementById("ParamIDSearch").style.display = "none";
  document.getElementById("IDSearch").style.display = "none";
  document.getElementById("DeleteID").style.display = "none";
}
function closeForm1() {
  document.getElementById("UpdStr").style.display = "none";
}

function openFormParamID() {
  document.getElementById("ParamIDSearch").style.display = "block";
  document.getElementById("AddNewStr").style.display = "none";
  document.getElementById("UpdStr").style.display = "none";
  document.getElementById("IDSearch").style.display = "none";
  document.getElementById("DeleteID").style.display = "none";
}
function closeFormParamID() {
  document.getElementById("ParamIDSearch").style.display = "none";
}
function openFormID() {
  document.getElementById("IDSearch").style.display = "block";
  document.getElementById("AddNewStr").style.display = "none";
  document.getElementById("UpdStr").style.display = "none";
  document.getElementById("ParamIDSearch").style.display = "none";
  document.getElementById("DeleteID").style.display = "none";
}
function closeFormID() {
  document.getElementById("IDSearch").style.display = "none";
}

function openFormDel() {
  document.getElementById("DeleteID").style.display = "block";
  document.getElementById("AddNewStr").style.display = "none";
  document.getElementById("UpdStr").style.display = "none";
  document.getElementById("ParamIDSearch").style.display = "none";
  document.getElementById("IDSearch").style.display = "none";
}
function closeFormDel() {
  document.getElementById("DeleteID").style.display = "none";
}
// function toggleForm(elementId) {
//   var div = document.querySelectorAll('.div');
  
//   div.forEach(function(div) {
//       div.style.display = (div.id === elementId) ? "block" : "none";
//   });
// }

// function closeForm(elementId) {
//   document.getElementById(elementId).style.display = "none";
// }


// function sortTable(columnIndex) {
//   var table, rows, switching, i, x, y, shouldSwitch;
//   table = document.getElementById("ScrollTable");
//   switching = true;

//   while (switching) {
//       switching = false;
//       rows = table.rows;

//       for (i = 1; i < rows.length - 1; i++) {
//           shouldSwitch = false;
//           x = rows[i].getElementsByTagName("td")[columnIndex];
//           y = rows[i + 1].getElementsByTagName("td")[columnIndex];

//           if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
//               shouldSwitch = true;
//               break;
//           }
//       }

//       if (shouldSwitch) {
//           rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
//           switching = true;
//       }
//   }
// }
//

//Реакт
const app = ReactDOMClient.createRoot((document.getElementById('app')))
app.render(<App />)
reportWebVitals();
//