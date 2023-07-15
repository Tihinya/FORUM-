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

const thread = document.getElementById("add-a-thread");
const details = document.getElementById("detailed-thread");

thread.addEventListener("click", function () {
  details.classList.toggle("show");
});

const threadTags = document.querySelectorAll(".tag-active");
let activeThread = threadTags[0];

threadTags.forEach(function (tag) {
  tag.addEventListener("click", function () {
    if (activeThread !== tag) {
      tag.classList.add("thread-subject-active");
      activeThread = tag;
    }
  });
});
