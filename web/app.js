"use strict";

// Theme toggle.
const themeToggle = document.getElementById("theme-toggle");
const savedTheme = localStorage.getItem("theme") || "dark";
document.documentElement.setAttribute("data-theme", savedTheme);
themeToggle.textContent = savedTheme === "dark" ? "☀️" : "🌙";

themeToggle.addEventListener("click", () => {
  const current = document.documentElement.getAttribute("data-theme");
  const next = current === "dark" ? "light" : "dark";
  document.documentElement.setAttribute("data-theme", next);
  localStorage.setItem("theme", next);
  themeToggle.textContent = next === "dark" ? "☀️" : "🌙";
});

const packSizesList = document.getElementById("pack-sizes-list");
const addPackBtn = document.getElementById("add-pack-btn");
const submitPacksBtn = document.getElementById("submit-packs-btn");
const calculateForm = document.getElementById("calculate-form");
const resultsTable = document.getElementById("results-table");
const resultsBody = document.getElementById("results-body");

// Render a single pack size input row.
function createPackRow(value) {
  const row = document.createElement("div");
  row.style.display = "flex";
  row.style.alignItems = "center";
  row.style.gap = "0.5rem";
  row.style.marginBottom = "0.5rem";

  const input = document.createElement("input");
  input.type = "number";
  input.min = "1";
  input.value = value || "";
  input.placeholder = "Pack size";
  input.required = true;
  input.style.marginBottom = "0";

  const removeBtn = document.createElement("button");
  removeBtn.type = "button";
  removeBtn.textContent = "x";
  removeBtn.className = "secondary outline";
  removeBtn.style.marginBottom = "0";
  removeBtn.style.width = "auto";
  removeBtn.style.padding = "0.5rem 0.75rem";
  removeBtn.addEventListener("click", () => row.remove());

  row.appendChild(input);
  row.appendChild(removeBtn);
  packSizesList.appendChild(row);
}

// Remove all child elements safely.
function clearChildren(element) {
  while (element.firstChild) {
    element.removeChild(element.firstChild);
  }
}

// Load current pack sizes from API.
async function loadPackSizes() {
  try {
    const res = await fetch("/api/packs");
    const data = await res.json();
    clearChildren(packSizesList);
    data.packs.forEach((size) => createPackRow(size));
  } catch (err) {
    console.error("Failed to load pack sizes:", err);
  }
}

// Add empty pack size row.
addPackBtn.addEventListener("click", () => createPackRow(""));

// Submit updated pack sizes.
submitPacksBtn.addEventListener("click", async () => {
  const inputs = packSizesList.querySelectorAll("input[type=number]");
  const packs = [];

  for (const input of inputs) {
    const val = parseInt(input.value, 10);
    if (!isNaN(val) && val > 0) {
      packs.push(val);
    }
  }

  if (packs.length === 0) {
    alert("At least one pack size is required. Changes were not saved.");
    await loadPackSizes();
    return;
  }

  try {
    const res = await fetch("/api/packs", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ packs }),
    });

    if (!res.ok) {
      const text = await res.text();
      alert("Error: " + text);
      return;
    }

    await loadPackSizes();
  } catch (err) {
    alert("Failed to update pack sizes: " + err.message);
  }
});

// Calculate packs for order.
calculateForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const order = parseInt(document.getElementById("order-items").value, 10);
  if (isNaN(order) || order <= 0) {
    alert("Enter a valid order quantity");
    return;
  }

  try {
    const res = await fetch("/api/calculate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ order }),
    });

    if (!res.ok) {
      const text = await res.text();
      alert("Error: " + text);
      return;
    }

    const data = await res.json();
    clearChildren(resultsBody);

    data.packs.forEach((entry) => {
      const row = document.createElement("tr");
      const packCell = document.createElement("td");
      packCell.textContent = entry.pack;
      const qtyCell = document.createElement("td");
      qtyCell.textContent = entry.quantity;
      row.appendChild(packCell);
      row.appendChild(qtyCell);
      resultsBody.appendChild(row);
    });

    resultsTable.style.display = "";
  } catch (err) {
    alert("Failed to calculate: " + err.message);
  }
});

// Load pack sizes on page load.
loadPackSizes();
