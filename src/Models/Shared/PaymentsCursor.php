<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class PaymentsCursor
{
    #[\JMS\Serializer\Annotation\SerializedName('cursor')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\PaymentsCursorCursor')]
    public PaymentsCursorCursor $cursor;

    public function __construct()
    {
        $this->cursor = new \formance\stack\Models\Shared\PaymentsCursorCursor();
    }
}