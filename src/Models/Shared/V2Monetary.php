<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2Monetary
{
    /**
     * The amount of the monetary value.
     *
     * @var int $amount
     */
    #[\JMS\Serializer\Annotation\SerializedName('amount')]
    #[\JMS\Serializer\Annotation\Type('int')]
    public int $amount;

    /**
     * The asset of the monetary value.
     *
     * @var string $asset
     */
    #[\JMS\Serializer\Annotation\SerializedName('asset')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $asset;

    public function __construct()
    {
        $this->amount = 0;
        $this->asset = '';
    }
}