---
title: Transactions with multiple sources
description: Sometimes you need to split a payment across multiple sources. Learn how to do this with Numscript.
---

import Prerequisites from '@site/docs/ledger/partials/numscript/_prerequisites.mdx';
import Prereqwarn from '@site/docs/ledger/partials/numscript/_prereq_warning.mdx';

Sometimes you need to split a payment from multiple sources. For example:

* At times there might not be enough money in an account you want to pay from. You want to specify a backup account.
* You need to spread costs across multiple accounts. For example, you might have a special marketing fund that partially covers certain payments.
* Your game players want to go in together on a shared purchase.

<Prerequisites />

## Basic transaction splitting

Since we're in the game land of Cones of Dunshire, let's consider that last case. Two players `donnameagle` and `tomhaverford` have decided [to treat themselves](https://www.youtube.com/watch?v=gSjM5B3QNlw) to a pony, which they will buy from us. The pony costs 75 coins, and they want to split the price evenly between themselves.

We can describe the transaction using Numscript. Create a file called `treat.num` with:

```numscript
send [COIN 75] (
  source = {
    50% from @player:donnameagle
    remaining from @player:tomhaverford
  }
  destination = @centralbank
)
```

And run it with

```shell
fctl ledger transactions num treat.num
```

<Prereqwarn />

Once you've run that transaction succesfully, let's have a look at `donnameagle`'s balance:


![`donnameagle` paid 38 and `tomhaverford` paid 37 coin](multi-source-1.png)

You should be able to see, as in the image above, that `donnameagle` paid 38 coins, and `tomhaverford` paid 37 coins.

### What's going on here?

First, we're telling Formance Ledger to take half of the 75 coin payment from `donnameagle`. Notice that 75, being an odd number, doesn't split evenly in half. Numscript is smart about this, and rounds amounts in a way that avoids rounding errors. Because `donnameagle` is listed first, they will pay the coin that remains after evenly subdividing: They will pay 38 coin.

:::info
[Floating point numbers are too imprecise for finance](https://www.youtube.com/watch?v=yZjCQ3T5yXo). Numscript avoids this problem by only using integer math for dividing payments up. The reference docs have [more detail on about Numscript's rounding algorithm](/ledger/reference/numscript/rounding/).
:::

Second, we're avoiding mistakes in our own calculations by telling Formance Ledger to pay whatever remains from `tomhaverford`, after `donnameagle` pays their share, by using the `remaining` keyword. Since `donnameagle` paid 38 coin, that leaves 37 coin for `tomhaverford` to pay.

## Specifying backup accounts

Splitting payments isn't the only use for specifying multiple source accounts. Sometimes you can't be sure that one account will have enough money to pay for a transaction, and you want to be able to specify a backup account to guarantee transaction success.

Suppose that `andydwyer` owes `aprilludgate` 100 coins for medical services after falling into a pit. But imagine moreover that `andydwyer` isn't sure that he has enough coin in his pocket, but he knows he has some stashed in a treasure chest. Here's how we can describe that transaction.

Create a file called `bills.num` with:

```numscript
send [COIN 100] (
  source = {
    @player:andydwyer
    @player:andydwyer:chest
  }
  destination = @player:aprilludgate
)
```

and run it with
```shell
fctl ledger transactions num bills.num
```

<Prereqwarn />

Run this script, and then [look at `andydwyer`'s balance](https://control.formance.com/accounts/player:andydwyer).

![`andydwyer` paid 50 from his pocket and 50 from his chest](multi-source-2.png)

You should be able to see, as in the image above, that 50 coin came from `andydwyer`'s primary account, and another 50 from the secondary account listed.

### What's going on here?

In this case, we have only specified a list of accounts as the source, without any apportionment. Formance Ledger will attempt to pay the entire amount from the first account listed. If there are insufficient funds in that account, then the remainder will be applied against the next account listed, and so on. You can list as many accounts as you like to draw funds from like this.

If the total in all the listed accounts remains insufficient to cover the movement, the transaction will fail and no money will be moved.

So, since there were only 50 coins in the first account, Formance Ledger looked at the second account, and took the remaining balance of 50 coins from it to complete the transaction.

## Nested sources

`donnameagle` and `tomhaverford` are not satisfied with just one pony—they want a second pony. But this time, `donnameagle` knows she has some upcoming transactions for which she needs to set aside some coin. She is only willing to pay 10 coin from her pocket. But she has some coin saved in a chest that she is willing to spend.

We could specify this new transaction as such.

```numscript
send [COIN 75] (
  source = {
    max [COIN 10] from @player:donnameagle
    37% from @player:donnameagle:chest
    remaining from @player:tomhaverford
  }
  destination = @centralbank
)
```

We use the keyword `max` to indicate that no more than the specified amount should be taken from `donnameagle`'s regular account.

But the remainder of the Numscript is annoyingly tedious: We had to compute that the share of the purchase price from `donnameagle`'s chest would be 37%. And if she doesn't have 10 coin in her pocket, that figure won't be correct.

There is a better way: Nested sources. Create a file called `pony2.num` with the following:

```numscript
send [COIN 75] (
  source = {
    50% from {
      max [COIN 10] from @player:donnameagle
      @player:donnameagle:chest
    }
    remaining from @player:tomhaverford
  }
  destination = @centralbank
)
```

and run it with
```shell
fctl ledger transactions num pony2.num
```

<Prereqwarn />

After running this script, and then [look at `donnameagle`'s balance](https://control.formance.com/accounts/player:donnameagle).

![`donnameagle` paid 10 from her pocket and 28 from her chest](multi-source-3.png)


### What's going on here?

Just like in the original pony purchase, we are asking Formance Ledger to split the cost 50-50 between `donnameagle` and `tomhaverford`. But, we've added a nested constraint on how we spend `donnameagle`'s coin. First, we use the `max` keyword to indicate that no more than 10 coins should be taken from her main account. Then we add a backup account—her treasure chest—to take the remainder of her share from.

## Going further

Numscript offers several different mechanisms for indicating how a transaction should be split amfrom different sources. This guide has just been a small taste of what's possible.

:::tip Dig deeper
Want to learn more about all the different ways to transact from multiple sources? [The reference docs](/ledger/reference/numscript/sources) have you covered!
:::
