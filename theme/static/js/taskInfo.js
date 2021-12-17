(() => {
  const rootGridElement = document.getElementById("root-grid");

  const toggleClass = (el, className) => {
    if (el.classList.contains(className)) {
      el.classList.remove(className);
    } else {
      el.classList.add(className);
    }
  };

  const removeClass = (className) => (el) => {
    el.classList.remove(className);
  };
  
  const addClickHandler = (el) => {
    el.addEventListener("click", () => {
      toggleClass(rootGridElement, "minimized");
      individualTaskElements.forEach(removeClass("active"));
      el.classList.add("active");
      el.scrollIntoView({behavior: "smooth", block: "start", inline: "nearest"});
    });
  };
  
  const individualTaskElements = document.querySelectorAll(".list-item");
  individualTaskElements.forEach(addClickHandler);
})();
