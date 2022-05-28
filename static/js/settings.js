let autoplaySelect = document.getElementById("autoplay");
let autoplayNextVidSelect = document.getElementById("autoplayNextVid");
let speedSelect = document.getElementById("speed");
let qualitySelect = document.getElementById("quality");

function loadSettings() {
  let plyrSettings = JSON.parse(localStorage.getItem("plyr")) || {};
  speedSelect.value = plyrSettings.speed || "1";
  qualitySelect.value = plyrSettings.quality || "0";
  autoplaySelect.checked = localStorage.getItem("autoplay") === "true" || false;
  autoplayNextVidSelect.checked = localStorage.getItem("autoplayNextVid") === "true" || false;
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

// SponsorBlock
let categories = ["sponsor", "selfpromo", "interaction", "intro", "outro", "preview", "filler"]

function updateSBSetting(category) {
  let categories = localStorage.getItem("sb_categories") || "";
  if (categories.includes(category)) {
    let re = new RegExp(`,?${category}`)
    localStorage.setItem("sb_categories", categories.replace(re, ""));
  } else if (categories.length == 0) {
    localStorage.setItem("sb_categories", categories + category);
  } else {
    localStorage.setItem("sb_categories", categories + "," + category);
  }

  let newCategories = localStorage.getItem("sb_categories")
  if (newCategories.startsWith(",")) {
    localStorage.setItem("sb_categories", newCategories.substring(1, 999));
  }
}

categories.forEach(category => {
  let select = document.getElementById(category)
  select.addEventListener("change", () => updateSBSetting(category))
  if (localStorage.getItem("sb_categories").includes(category)) {
    select.checked = true
  }
})