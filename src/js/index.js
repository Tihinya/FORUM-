const showPopupButton = document.getElementById("show-button");
const popup = document.getElementById("popup");

showPopupButton.addEventListener("click", function () {
  popup.classList.toggle("show");
});

const notifyButton = document.getElementById("notification-button");
const notifyPopup = document.getElementById("notify-popup");

notifyButton.addEventListener("click", function () {
  notifyPopup.classList.toggle("show");
});

var subjects = document.querySelectorAll(".nav__options_hover");
var activeSubj = subjects[0];

subjects.forEach(function (subject) {
  subject.addEventListener("click", function () {
    if (activeSubj !== subject) {
      activeSubj.classList.remove("nav__options_active");
      subject.classList.add("nav__options_active");

      activeSubj = subject;
    }
  });
});
