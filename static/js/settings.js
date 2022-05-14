let autoplaySelect = document.getElementById("autoplay");
let autoplayNextVidSelect = document.getElementById("autoplayNextVid");
let speedSelect = document.getElementById("speed");
let qualitySelect = document.getElementById("quality");

function loadSettings() {
  let plyrSettings = JSON.parse(localStorage.getItem("plyr")) || {};
  speedSelect.value = plyrSettings.speed || "1";
  qualitySelect.value = plyrSettings.quality || "0";
  autoplaySelect.checked = localStorage.getItem("autoplay") || false;
  autoplayNextVidSelect.checked = localStorage.getItem("autoplayNextVid") || false;
}
loadSettings()

autoplaySelect.addEventListener("change", () => {
  localStorage.setItem("autoplay", autoplaySelect.checked);
})
autoplayNextVidSelect.addEventListener("change", () => {
  localStorage.setItem("autoplayNextVid", autoplayNextVidSelect.checked);
})
qualitySelect.addEventListener("change", () => {
  let plyrStorage = localStorage.getItem("plyr");
  let data = plyrStorage ? JSON.parse(plyrStorage) : {};
  data.quality = qualitySelect.value * 1;
  localStorage.setItem("plyr", JSON.stringify(data));
})
speedSelect.addEventListener("change", () => {
  let plyrStorage = localStorage.getItem("plyr");
  let data = plyrStorage ? JSON.parse(plyrStorage) : {};
  data.speed = speedSelect.value * 1;
  localStorage.setItem("plyr", JSON.stringify(data));
})