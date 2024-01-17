//отображение формы
function toggleForm(formId) {
  var allForms = document.querySelectorAll('.form');
  allForms.forEach(function(form) {
      form.style.display = 'none';
  });
  var selectedForm = document.getElementById(formId);
  if (selectedForm) {
      selectedForm.style.display = 'block';
  }
}

//Кнопка закрытия формы
function CloseForm(formId) {
  var selectedForm = document.getElementById(formId);
  if (selectedForm) {
    selectedForm.style.display = 'none';
  }
}
