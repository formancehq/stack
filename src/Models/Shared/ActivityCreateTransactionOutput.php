<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class ActivityCreateTransactionOutput
{
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\OrchestrationTransaction')]
    public OrchestrationTransaction $data;
    
	public function __construct()
	{
		$this->data = new \formance\stack\Models\Shared\OrchestrationTransaction();
	}
}