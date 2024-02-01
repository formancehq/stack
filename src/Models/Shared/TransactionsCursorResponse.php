<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TransactionsCursorResponse
{
	#[\JMS\Serializer\Annotation\SerializedName('cursor')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\TransactionsCursorResponseCursor')]
    public TransactionsCursorResponseCursor $cursor;
    
	public function __construct()
	{
		$this->cursor = new \formance\stack\Models\Shared\TransactionsCursorResponseCursor();
	}
}