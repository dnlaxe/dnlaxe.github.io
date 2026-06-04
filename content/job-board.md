---
title: "Job Board"
date: "2026-06-01"
type: "project"
slug: "project-job-board"
tags: ["node", "express", "handlebars"]
---

[_korea-jobs-board.vercel.app_](https://korea-jobs-board.vercel.app/jobs/board)

# Intention

I wanted to create a simple job board. The main goal was for it to be frictionless and easy to post and manage posts. In order to achieve this, users do not have to create an account. With an email alone, they are able to submit posts. The email is only needed so that the user can recieve an email in which there is a link through which they can manage their post. I also decided to use a multistep form as a quick way to construct posts.

# Stack

I chose to use Node/Express with Typescript for the backend. I wanted the page to server side rendered to help with SEO.

The muiltstep form was made with vanilla javascript.

# Key dependencies

- Tailwind for styling.
- Drizzle and Postgres for database.
- Zod for validation.
- Handlebars for templating.
- Pino for logging.
- Compression, cookie-parser, rate-limiter-flexible, helmet as recommended by Express docs.
- Resend for email.

# Flags

I added flags in development so that mock payment and mock email can be toggled.

# Error Handling

I decided to try and implement an error handling system using discriminated unions.

```
export type Result<T> =
  | { success: true; data: T }
  | { success: false; error: appError };
```

Through using discriminated unions, TypeScript will throw an error during development if accessing `.data` unless `.success` has been checked first first.

Rejected promises in async functions should be caught by a `try/catch`. I wanted to turn these rejections into a standard return values.

Instead of throwing errors, errors are caught and consumed within the function. As far as the app is concerned, the promise was resolved successfully. This error object can then be dealt with manually in an explicit way.

I use `errorHandler`, `notFoundHandler` and `globalErrorHandler` middleware to catch any unforeseen errors.

# Architecture

Routes → http only, middleware and delegates to controller.

Controllers → http only, reads req, calls services, then writes res.

Services → no http, only business logic, returns result objects

Repos → only db queries

I decided to group them together by feature (jobs, admin etc). Apart from for the repos which has its own folder.

# Magic link system

# Multistep form

The form is handled by `form.js`. Whenever the user interacts with the form, a state object is updated. This is the source of frontend truth.

Each step in the form is a choice of radio buttons. In two cases the choice of radio button informs the options shown in the next step. This could easily get messy so I used a function to track form's progress. This function is a check so it returns a boolean.

```
function isStepComplete() {
  const steps = stepsConfig[current]; // Get the fields in that step

  // Conditional logic for special cases (in this case it's the contact step)
  if (current === STEPS.CONTACT) {
    if (state.contactMethod === "relay") return true; // No extra fields needed
    if (state.contactMethod === "link") return state.contactUrl !== ""; // Must have a URL
    return false;
  }

  // Check each field in that step has been selected and is valid
  return (steps.fields ?? []).every((field) => {
    const f = state[field];
    return f !== null && f !== "" && f !== undefined;
  });
}
```

After user has been through each step, they are shown a review (`buildReviewPage`) of their selections and they are able to go back to a specific step to make changes. This means they can reuse the existing steps however there will be certain changes as they are now in review mode. They should only be able to switch between the review and a certain step. By using `isEdit` the form's UI logic,routing and available buttons, can be changed. This also means taht if a parent option is changed, the review will be able to tell the user to update a child option.

If the server returns an error after attmepted submission, the page will be reloaded directly into review mode (`isEdit` = true) and the form will be rehydrated with error feedback for the user (`rebuildState`).

Creating this form was interesting however I would like to remake it again in a simpler, easier to read way.

# Audit logging

I wanted a simple way of keeping a record of certain events. In an audit event table in the database, I am able to track these events with details. These events are then shown in the admin dashboard.

It is particularly useful when with the `expireOverduePosts` function. When the job board is loaded, there is a quick check to see if any job posts have expired. If they are, they are marked as expired in the database so that they no longer appear. By recording this in audit events, admin is able to know that this has happened and when.

# View tracking

A view count every time a job is viewed happens as an example of collecting statistics.

# Mock payment and mock email

I set up a mock payment and mock email flow that can be toggled. Resend is responsible for the real emails.

# Database and cache

I use a Postgres database with drizzle.

I set up a cache to reduce the number of times the database is called. A simple check `if (Date.now() > activePostsCache.expiresAt)` decides whether to serve the cache or empty and call the database for the latest data.

```
const TIME_TO_EXPIRY = 60_000;

interface LivePostsCache {
  posts: LivePostRow[];
  expiresAt: number;
}
```

# Environment variable validation

Zod validates `process.env` on startup and the application will not start if the environment isn't configured correctly. This reduces the chance of environment variables being the cause of a bug.

These variables are then extracted into a config object to passed around the project.

# Config flags

This config object can then be used to create useful flags.

```
export const isDevelopment = config.node_env === "development";
export const isProduction = config.node_env === "production";
export const isBasicAuthEnabled = config.basic_auth === true;
```

# Middleware

# Basic Auth

# Templating

Handlebars are used for templating. Repeated components are extracted into partials. Helpers, e.g. `json: (obj) => JSON.stringify(obj)` are added to add more functionality in html.

# Health checks

`showReadiness` checks the database is still connected. `showLiveness` shows process is alive.

# Startup

The server doesn't start unless the database is connected and the enviroment variables are valid.

# Graceful shutdown

`shutdownHandler` attempts to shut the app down in orderly way if an unexpected error occurs. By using an isShuttingDown flag, the readiness probe will start returning `503 Service Unavailable` which will tell the load balancer to stop sending new requests.

`shutdown.ts` then drops idle connections, waits for active requests to finish, and closes connections (the database).

Adding a `SHUTDOWN_TIMEOUT_MS` timer ensures the shutdown doesn't hang ane will terminate hanging processes.

# Zod

Zod validates any data coming into the backend from the client.

# Session

I wanted users to be able to add multiple posts and have them persist between refreshes similar to a shopping cart on commerce website. After receiving and validating the user's email (`/start`), a session is created and the user can start creating (`/form`) and saving posts. The form cannot be accessed until an email has been submitted in '/start'.

# Vercel hosting, postgres, vercel.json

# Summary

If I were to make something like this again I would like to:

- Use dependency injection
- Use deep modules for the architecture
- Use less dependencies
- Create a simple design system beforehand that could then easily updated later
- More structured frontend
