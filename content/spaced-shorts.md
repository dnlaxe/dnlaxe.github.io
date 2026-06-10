---
title: "Spaced Shorts"
date: "LAST UPDATED: 2026-06-09"
type: "project"
slug: "project-spaced-shorts"
tags: ["react", "language learning"]
---

## Live Demo

[_spaced-shorts.vercel.app_](https://spaced-shorts.vercel.app/)

## Intention

`Comprehensible Input` is a language acquisition theory by linguist Stephen Krashen. It suggests we learn languages best not by memorizing grammar rules, but by understanding messages. The core idea is to consume listening and reading content you mostly understand, containing just a few new words or structures to stretch your current proficiency.

I decided to take this idea and combine it with an Anki-style flashcard app and Youtube shorts. Youtube shorts are inherently addictive and easy to watch, offering an easy bitesized way to consume input. Then by using an Anki-style algorithm to track progress, shorts can be reshown systematically to reinforce the input.

## Stack

As this is a frontend project, I decided to use React, React Router and Vite, along with TypeScript for type safety and Tailwind CSS for styling. I initially added an Express backend for persistence, but later chose to simplify the project by storing data in local storage through a custom hook.

## Architecture

The application is split into three responsibilities:

- Playlist management.
- Session construction.
- Spaced repetition and queue management.

`buildSessionShorts` determines which shorts are due when a user starts a practice session. `useQueueHandler` manages the active queue and updates short metadata in response to user feedback. Playlist data is persisted in local storage through a custom hook.

## Spaced Repetition Algorithm

Every short contains information.

```
export type Short = {
  id: string;
  due: number;
  intervals: number;
  ease: number;
  state: string;
  stepIndex: number;
  reps: number;
  lapses: number;
};
```

`id` uniquely identifies the short.

`due` a number (timestamp) of when the short should be shown again.

`intervals` is the amount of days between reviews. The intervals between reviews should grow. The length of intervals change depending on the user's response combined with the `ease` factor.

`ease` a fixed number which is used to calculate interval length. In Anki, this value changes, but in this project it is kept fixed for simplicity.

_Ease and intervals are used only for shorts in review._

`state` a short's state can be new, learning, review, relearning. Each state has different rules.
State changes:

```
new
 ↓
learning
 ↓
review
 ↑
relearning
```

New shorts change to learning when first shown.
Learning changes to review when `stepIndex >= 2`
Review changes to relearning if user's response is `again`
Relearning changes back to review for any user response other than `again`.

`stepIndex` each short can have a stepIndex of 0-2. This number is changed by user's response. This is used to change state.

- `Again` stepIndex = 0

- `Hard` no change

- `Medium` stepIndex += 1

- `Easy` stepIndex += 2

_stepIndex is only used for learning/relearning shorts_

`reps` counts the number of times a short has been reviewed.

`lapses`counts the number of times a short has been failed.

_reps and lapses are currently unused._

# Learning Phase

Every new short starts with step: 0 and state: learning. This changes depending on user's feedback:

`Again` step = 0, lapses++

`Hard` no change

`Medium` step increments by 1

`Easy` step increments by 2

When `step >= 2`, the short graduates to the `review` state.

# Review Phase

shorts that complete learning enter the `review` state. User feedback changes the review interval:

`Again` state = relearning, step = 0, lapses++, interval = 0

`Hard` interval × 1.2

`Medium` interval × ease

`Easy` interval × (ease + 0.15)

After a successful review (`Hard`, `Medium`, or `Easy`), the short's due date is updated based on the new interval.

# Relearning Phase

When a short in the `review` state receives an `Again` rating, it enters the `relearning` state with `step = 0`.

Shorts in the `relearning` state remain in the current session until they successfully return to the `review` state.

## Watch

When the user clicks practice they are taken to page to see the shorts with 5 buttons (exit, again, hard, okay, easy) below the short.

In the background, a list of session cards is constructed. Learning and relearning shorts are always included. Review shorts are included only if their due timestamp is earlier than the current time, and are sorted by due date. New shorts are added afterwards. User-configurable limits determine how many review and new shorts are included in a session.

The cards are managed by a custom hook, `useQueueHandler`. User responses mutate the metadata of the current short, and the resulting state determines how the queue is organised. Shorts in the learning and relearning states are appended back onto the queue, while shorts in the review state leave the session and are scheduled for a future date. The hook accumulates all mutated shorts and, once the session is complete, returns a merged collection that can be persisted to local storage or a database.

## Playlist management

Users are able to create playlists and then add shorts by clicking add and then copy and pasting in the full youtube short url address. The final part of this url is parsed and kept, this can be then used later to rebuild the address when showing the short.

Shorts can be added and deleted. Whole playlists can be deleted. Settings allow a user to choose limits on how many new cards and review cards are shown to them a day for each playlist.

## Future

If I were to remake this again, I would:

- Separate the scheduling algorithm from the queue management logic so that the two concerns are independent.
- Add a way to add a whole playlist of shorts.
- Add a page showing statistics.
