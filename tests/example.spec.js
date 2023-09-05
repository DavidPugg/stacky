// @ts-check
const { test, expect, chromium } = require("@playwright/test");

const pageURL = "http://localhost:3000";
const mainPageURL = `${pageURL}/`;
const discoverPageURL = `${pageURL}/discover`;
const createPageURL = `${pageURL}/create`;
const registerPageURL = `${pageURL}/register`;
const loginPageURL = `${pageURL}/login`;

let browser;
let authContext;

test.beforeAll(async () => {
  browser = await chromium.launch();
  authContext = await browser.newContext();

  await register(authContext);
  await login(authContext);
});

test.afterAll(async () => {
  await browser.close();
});

//Before authentication

test("should not be able to go to certain routes", async ({ page }) => {
  await page.goto(mainPageURL);
  expect(page.url()).toEqual(discoverPageURL);

  await page.goto(createPageURL);
  expect(page.url()).toEqual(discoverPageURL);
});

test("should show sign in link", async ({ page }) => {
  await page.goto(mainPageURL);
  await page.click("a:has-text('Sign in')");
  await page.waitForURL(loginPageURL);
  expect(page.url()).toEqual(loginPageURL);
});

test("should not show avatar", async ({ page }) => {
  await page.goto(mainPageURL);
  const avatarLocator = page.getByTestId("avatar");
  const avatarCount = await avatarLocator.count();
  expect(avatarCount).toBe(0);
});

test("register", async ({ context }) => {
  return await register(context);
});

test("login", async ({ context }) => {
  return await login(context);
});

//After authentication

test("has title", async () => {
  const page = await authContext.newPage();
  await page.goto(mainPageURL);
  await expect(page).toHaveTitle(/Stacky/);
});

test("should show no posts placeholder", async () => {
  const page = await authContext.newPage();
  await page.goto(mainPageURL);
  expect(page.getByText("No posts yet")).toBeDefined();
});

test("should not show signin", async () => {
  const page = await authContext.newPage();
  await page.goto(mainPageURL);
  const signInText = await page.$("text=Sign in");
  expect(signInText).toBeNull();
});

test("should show avatar", async () => {
  const page = await authContext.newPage();
  await page.goto(mainPageURL);
  const avatarLocator = page.getByTestId("avatar");
  const avatarCount = await avatarLocator.count();
  expect(avatarCount).toBe(1);
});

//Auth functions

async function register(context) {
  const page = await context.newPage();

  let registered = false;

  await page.goto(registerPageURL);
  await page.fill("#username", "playwrightTest");
  await page.fill("#email", "playwright@playwright.com");
  await page.fill("#password", "12345678");
  await page.click("button:has-text('Register')");

  if (
    page.getByText("Username already in use") ||
    page.getByText("Email already in use")
  ) {
    await page.click("a:has-text('Login')");
  } else {
    registered = true;
  }

  await page.waitForURL(loginPageURL);

  if (registered) {
    expect(page.getByText("User created")).toBeDefined();
  }

  expect(page.url()).toBe(loginPageURL);
}

async function login(context) {
  const page = await context.newPage();

  await page.goto(loginPageURL);
  await page.fill("#username", "playwrightTest");
  await page.fill("#password", "12345678");
  await page.click("button:has-text('Login')");
  await page.waitForURL(mainPageURL);
  expect(page.getByText("Successfully logged in")).toBeDefined();
  expect(page.url()).toBe(mainPageURL);
  const cookies = await page.context().cookies();
  const session_id = cookies.find((cookie) => cookie.name === "session_id");
  if (!session_id) return;
  context.addCookies([session_id]);
}
