function btn1() {
  var tab = document.getElementById("div_tab")
  var btn = document.getElementById("btn1")
  btn.classList.toggle('btn_op')
  btn.classList.toggle('btn_close')
  tab.classList.toggle('table_div_vis')
  tab.classList.toggle('table_div_invis')
  }

  function callDatabase() {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:3000", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log("good")
        }
    };

    xhr.send("action=database()");
}
