<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2LedgerInfoResponse
{
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\V2LedgerInfo')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?V2LedgerInfo $data = null;
    
	public function __construct()
	{
		$this->data = null;
	}
}
