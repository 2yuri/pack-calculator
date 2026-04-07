import { test, expect } from "@playwright/test";

test("create product, add packs, and calculate order", async ({ page }) => {
  const productName = `E2E Product ${Date.now()}`;

  // Create product
  await page.goto("/products");
  await page.getByRole("button", { name: "New Product" }).click();
  await page.locator("#product-name").fill(productName);
  await page.getByRole("button", { name: "Create" }).click();
  await expect(page.getByText("Product created")).toBeVisible();
  await expect(page.getByRole("cell", { name: productName })).toBeVisible();

  // Create Packs
  await page.getByRole("row", { name: productName }).click();
  await expect(page.getByRole("heading", { name: /Packs/ })).toBeVisible();
  const productId = page.url().match(/\/products\/(\d+)\/packs/)![1];

  for (const size of [250, 500, 1000]) {
    await page.getByRole("button", { name: "Add Pack" }).first().click();
    await page.locator("#pack-size").fill(String(size));
    await page
      .getByRole("dialog")
      .getByRole("button", { name: "Add Pack" })
      .click();
    await expect(page.getByRole("dialog")).toBeHidden();
  }

  // Calculate
  await page.getByRole("link", { name: /order/i }).click();
  await expect(page.getByText("Order Calculator")).toBeVisible();

  const productCard = page.locator('[data-slot="card"]', {
    hasText: `#${productId}`,
  });
  await productCard.locator('input[type="number"]').fill("263");

  await page.getByRole("button", { name: /calculate/i }).click();

  // Verify results
  await expect(page.getByRole("heading", { name: "Results" })).toBeVisible();

  const resultCard = page.locator('[data-slot="card"]', {
    has: page.locator('[data-slot="card-header"]', { hasText: productName }),
  });
  await expect(resultCard).toBeVisible();
  await expect(resultCard.getByText("263")).toBeVisible();
  await expect(resultCard.getByRole("table")).toBeVisible();

  const packRows = resultCard.getByRole("row").filter({ hasText: /x\d+/ });
  await expect(packRows.first()).toBeVisible();
});
