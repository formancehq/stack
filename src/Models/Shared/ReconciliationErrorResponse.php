<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class ReconciliationErrorResponse
{
	#[\JMS\Serializer\Annotation\SerializedName('details')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $details = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('errorCode')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $errorCode;
    
	#[\JMS\Serializer\Annotation\SerializedName('errorMessage')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $errorMessage;
    
	public function __construct()
	{
		$this->details = null;
		$this->errorCode = "";
		$this->errorMessage = "";
	}
}
