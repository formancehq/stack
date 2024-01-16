---
title: Transactions with multiple destinations
description: Sometimes you need to split a payment across multiple destinations. Learn how to do this with Numscript.
---

import Prerequisites from '@site/docs/ledger/partials/numscript/_prerequisites.mdx';
import Prereqwarn from '@site/docs/ledger/partials/numscript/_prereq_warning.mdx';


Sometimes you need to split a payment across multiple destinations. For example:

* You charged VAT on an order, and so you need to split off a percentage of a payment into a dedicated VAT account.
* You have a marketplace, and a customer has made a single purchase of items from multiple vendors.
* Your game players need to split the loot from a group raid.

<Prerequisites />

## Basic transaction splitting

Since we're in the game land of Cones Dunshire, let's consider that last case. As the ledgerman, we want a map of a coastal region, and we're willing to pay 75 coins to someone to get it done. `leslieknope` and `annperkins`, a group of two surveyors working together, rise to the task, and create that map. Now, we want to split the reward between them.

We can describe the transaction using Numscript. Create a file called `split.num` with:

```numscript
send [COIN 75] (
  source = @centralbank
  destination = {
    50% to @player:leslieknope
    remaining to @player:annperkins
  }
)
```

And run it with

```shell
fctl ledger transactions num split.num
```

<Prereqwarn />

Once you've run that transaction succesfully, let's have a look at `leslieknope`'s balance:


![`leslieknope` gets 38 and `annperkins` gets 37 coin](multi-destination-1.png)

You should be able to see, as in the image above, that `leslieknope` received 38 coins, and `annperkins` received 37 coins.

### What's going on here?

First, we're sending half of the 75 coin payment to `leslieknope`. Notice that 75, being an odd number, doesn't split evenly in half. Numscript is smart about this, and rounds amounts in a way that avoids rounding errors. Because `leslieknope` is listed first, they will get the coin that remains after evenly subdividing: They will receive 38 coin.

:::info Formance Ledger uses integer math
[Floating point numbers are too imprecise for finance](https://www.youtube.com/watch?v=yZjCQ3T5yXo). Numscript avoids this problem by only using integer math for dividing payments up. The reference docs have [more detail on about Numscript's rounding algorithm](/ledger/reference/numscript/rounding/).
:::

Second, we're avoiding mistakes in our own calculations by telling Numscript to send whatever remains to `annperkins`, after `leslieknope` gets their share, by using the `remaining` keyword. Since `leslieknope` received 38 coin, that leaves 37 coin to distribute to `annperkins`.

## Nested transaction splitting.

Let's take the previous scenario, and add a twist. Let's suppose we need to withhold taxes from the payment—Dunshire imposes a flat 15% sales tax for goods and services. We could modify `split.num` to reflect a three-way transaction:

```numscript
send [COIN 75] (
  source = @centralbank
  destination = {
    15% to @salestax
    43% to @player:leslieknope
    remaining to @player:annperkins
  }
)
```

But there is a problem: Scripting the transaction this way requires us to manually do the math to figure out that `leslieknope` should get 42.5% of the transaction (which we've had to round up).

There's a better way: Nested destinations. Using nested destinations allows us to be clearer about our intent, and leaves much of the math off to Numscript to sort out:

```numscript
send [COIN 75] (
  source = @centralbank
  destination = {
    15% to @salestax
    remaining to {
        50% to @player:leslieknope
        remaining to @player:annperkins
    }
  }
)
```

<Prereqwarn />

Run this script, and then [look at `leslieknope`'s balance](https://control.formance.com/accounts/player:leslieknope).

![`leslieknope` gets 38 and `annperkins` gets 37 coin](multi-destination-2.png)

You can see that Numscript has worked out all of necessary calculations so that the entire payment is allocated correctly to all destinations, with no remainder.

## Going further

Numscript offers several different mechanisms for indicating how a transaction should be split among different destinations, this guide has just been a small taste of what's possible.

:::tip Dig deeper
Want to learn more about all the different ways to split a transaction? [The reference docs](/ledger/reference/numscript/destinations) have you covered!
:::
