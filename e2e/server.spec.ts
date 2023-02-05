import { test, expect, request } from "@playwright/test";

const testUrl = (path: string): string => {
  return `http://test.localhost:8080${path}`;
};

test("test get request", async ({ page }) => {
  const response = await page.goto(testUrl("/"));

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("test post request with valid json data", async ({ request }) => {
  const response = await request.post(testUrl("/"), {
    data: {
      username: "test",
      password: "test",
    },
    headers: {
      "Content-Type": "application/json",
    },
  });

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("test post request with valid form data", async ({ request }) => {
  const response = await request.post(testUrl("/"), {
    data: {
      username: "test",
      password: "test",
    },
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
  });

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("test put request", async ({ request }) => {
  const response = await request.put(testUrl("/"));

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("test patch request", async ({ request }) => {
  const response = await request.patch(testUrl("/"));

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("test delete request", async ({ request }) => {
  const response = await request.delete(testUrl("/"));

  expect(response?.ok());

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("302 redirect request", async ({ page }) => {
  const response = await page.goto(testUrl("/redirect-302"));

  expect(response?.status()).toEqual(200);

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});

test("307 redirect request", async ({ page }) => {
  const response = await page.goto(testUrl("/redirect-307"));

  expect(response?.status()).toEqual(200);

  const resJson = await response?.json();

  expect(resJson).toEqual({ message: "ok" });
});
