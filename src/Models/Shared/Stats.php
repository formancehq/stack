<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Stats
{
    /**
     *
     * @var int $accounts
     */
    #[\JMS\Serializer\Annotation\SerializedName('accounts')]
    public int $accounts;

    /**
     *
     * @var int $transactions
     */
    #[\JMS\Serializer\Annotation\SerializedName('transactions')]
    public int $transactions;

    /**
     * @param  ?int  $accounts
     * @param  ?int  $transactions
     */
    public function __construct(?int $accounts = null, ?int $transactions = null)
    {
        $this->accounts = $accounts;
        $this->transactions = $transactions;
    }
}