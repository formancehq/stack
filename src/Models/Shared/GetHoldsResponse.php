<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class GetHoldsResponse
{
    /**
     *
     * @var GetHoldsResponseCursor $cursor
     */
    #[\JMS\Serializer\Annotation\SerializedName('cursor')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\GetHoldsResponseCursor')]
    public GetHoldsResponseCursor $cursor;

    /**
     * @param  ?GetHoldsResponseCursor  $cursor
     */
    public function __construct(?GetHoldsResponseCursor $cursor = null)
    {
        $this->cursor = $cursor;
    }
}