function init() {
  const colors = [
    "#1dbf73",
    "#80489C",
    "#432C7A",
    "#FF8FB1",
    "#FD841F",
    "#B2A4FF",
    "#3F4E4F",
    "#2C3639",
  ];
  document.body.style.backgroundColor =
    colors[Math.floor(Math.random() * colors.length)];
}

window.onload = init();
