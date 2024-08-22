<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class PoolBalances
{
    /**
     * $balances
     *
     * @var array<PoolBalance> $balances
     */
    #[\JMS\Serializer\Annotation\SerializedName('balances')]
    #[\JMS\Serializer\Annotation\Type('array<\formance\stack\Models\Shared\PoolBalance>')]
    public array $balances;

    /**
     * @param  ?array<PoolBalance>  $balances
     */
    public function __construct(?array $balances = null)
    {
        $this->balances = $balances;
    }
}