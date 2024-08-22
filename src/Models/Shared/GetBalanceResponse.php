<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class GetBalanceResponse
{
    /**
     *
     * @var BalanceWithAssets $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\BalanceWithAssets')]
    public BalanceWithAssets $data;

    /**
     * @param  ?BalanceWithAssets  $data
     */
    public function __construct(?BalanceWithAssets $data = null)
    {
        $this->data = $data;
    }
}