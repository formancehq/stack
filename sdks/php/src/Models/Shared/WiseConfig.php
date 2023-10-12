<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class WiseConfig
{
	#[\JMS\Serializer\Annotation\SerializedName('apiKey')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $apiKey;
    
    /**
     * The frequency at which the connector will try to fetch new BalanceTransaction objects from Wise API.
     * 
     * 
     * 
     * @var ?string $pollingPeriod
     */
	#[\JMS\Serializer\Annotation\SerializedName('pollingPeriod')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $pollingPeriod = null;
    
	public function __construct()
	{
		$this->apiKey = "";
		$this->pollingPeriod = null;
	}
}
