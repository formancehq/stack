<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2CreateTransactionResponse
{
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\V2Transaction')]
    public V2Transaction $data;

    public function __construct()
    {
        $this->data = new \formance\stack\Models\Shared\V2Transaction();
    }
}