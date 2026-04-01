"use strict";

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
  // TODO: implement when PUT /api/packs is ready.
});

// Calculate packs for order.
calculateForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  // TODO: implement when POST /api/calculate is ready.
});

// Load pack sizes on page load.
loadPackSizes();
