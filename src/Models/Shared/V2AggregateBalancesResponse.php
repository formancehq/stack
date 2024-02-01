<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2AggregateBalancesResponse
{
    /**
     * $data
     * 
     * @var array<string, int> $data
     */
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<string, int>')]
    public array $data;
    
	public function __construct()
	{
		$this->data = [];
	}
}