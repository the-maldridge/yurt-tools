(() => {
  let activeTaskElement = document.querySelector(".active");
  const rootDetailElement = document.getElementById("root-detail");
  const rootGridElement = document.getElementById("root-grid");
  const individualTaskElements = document.querySelectorAll(".list-item");

  const addClass = (el, className) => {
    if (!el.classList.contains(className)) {
      el.classList.add(className);
    }
  };

  const resetActiveClasses = () =>
    individualTaskElements.forEach((el) => el.classList.remove("active"));

  /*
   * This handler does a few things:
   * 1. Makes the root grid work as a one column grid, ie scrollable list
   * 2. Add an active class to only the clicked element
   * 3. Scrolls the element into view
   */
  const selectTask = (el) => {
    activeTaskElement = el;
    addClass(rootGridElement, "minimized");
    addClass(rootDetailElement, "maximized");
    addClass(el, "active");
    el.scrollIntoView({
      behavior: "smooth",
      block: "start",
      inline: "nearest",
    });
    const { namespace, job, group, task } = el.dataset;
    history.pushState(
      null,
      document.title,
      `/taskinfo/view/${namespace}/${job}/${group}/${task}`
    );
    fetch(`/taskinfo/view/${namespace}/${job}/${group}/${task}/details`)
      .then((response) => response.text())
      .then((responseHTML) => rootDetailElement.innerHTML = responseHTML)
  };

  const unselectTask = () => {
    activeTaskElement = null;
    rootGridElement.classList.remove("minimized");
    rootDetailElement.classList.remove("maximized");
    rootDetailElement.innerHTML = "";
    history.pushState(null, document.title, "/taskinfo/view");
  };

  const addClickHandler = (el) => {
    el.addEventListener("click", () => {
      resetActiveClasses();
      if (activeTaskElement === el) {
        return unselectTask();
      }
      selectTask(el);
    });
  };

  individualTaskElements.forEach(addClickHandler);
})();
