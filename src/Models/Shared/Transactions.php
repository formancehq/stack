<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Transactions
{
    /**
     * $transactions
     *
     * @var array<TransactionData> $transactions
     */
    #[\JMS\Serializer\Annotation\SerializedName('transactions')]
    #[\JMS\Serializer\Annotation\Type('array<\formance\stack\Models\Shared\TransactionData>')]
    public array $transactions;

    /**
     * @param  ?array<TransactionData>  $transactions
     */
    public function __construct(?array $transactions = null)
    {
        $this->transactions = $transactions;
    }
}