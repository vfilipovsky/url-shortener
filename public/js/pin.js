const codes = document.querySelectorAll(".code");
var pin = [];

codes[0].focus();

codes.forEach((code, idx) => {
  code.addEventListener("keydown", async (e) => {
    if (e.key >= 0 && e.key <= 9) {
      pin[idx] = e.key;
      codes[idx].value = "";
      if (idx < 3) {
        setTimeout(() => codes[idx + 1].focus(), 10);
      } else {
        setTimeout(() => document.activeElement.blur(), 10);
        await check();
      }
    } else if (e.key === "Backspace") {
      if (idx > 0) {
        setTimeout(() => codes[idx - 1].focus(), 10);
      }
    } else {
      codes[idx].value = "";
    }
  });
});

async function check() {
  const requestUrl = "{{.Url}}/api/v1/url/{{.Code}}";

  const res = await fetch(requestUrl, {
    method: "POST",
    body: JSON.stringify({
      pin: pin.join(""),
    }),
  });

  if (res.status !== 200) {
    const info = document.getElementById("info");
    info.style.visibility = "";
    info.style.color = "#DD5353";
    reset();
    return;
  }

  const url = await res.json();

  window.location.replace(url);
}

function reset() {
  for (let i = 0; i < codes.length; i++) {
    codes[i].value = "";
  }

  codes[0].focus();
}
