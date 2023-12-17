function btn1() {
  var tab = document.getElementById("div_tab")
  var btn = document.getElementById("btn1")
  btn.classList.toggle('btn_op')
  btn.classList.toggle('btn_close')
  tab.classList.toggle('table_div_vis')
  tab.classList.toggle('table_div_invis')
  }
function openForm() {
  document.getElementById("myForm").style.display = "block";
}

function closeForm() {
  document.getElementById("myForm").style.display = "none";
}
