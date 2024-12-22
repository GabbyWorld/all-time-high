export const designWidth = 375
export const rootFontSize = 16

function setRemUnit() {
  const docEl = document.documentElement;
  const clientWidth = docEl.clientWidth || window.innerWidth;
  const rem = (clientWidth / designWidth) * rootFontSize;
  docEl.style.fontSize = rem + "px";
}

setRemUnit();

window.addEventListener("resize", setRemUnit);
window.addEventListener("pageshow", (e: PageTransitionEvent) => {
  if (e.persisted) {
    setRemUnit();
  }
});
