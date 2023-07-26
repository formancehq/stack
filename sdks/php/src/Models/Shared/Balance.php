<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Balance
{
	#[\JMS\Serializer\Annotation\SerializedName('accountId')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $accountId;
    
	#[\JMS\Serializer\Annotation\SerializedName('balance')]
    #[\JMS\Serializer\Annotation\Type('int')]
    public int $balance;
    
	#[\JMS\Serializer\Annotation\SerializedName('createdAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    public \DateTime $createdAt;
    
	#[\JMS\Serializer\Annotation\SerializedName('currency')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $currency;
    
	#[\JMS\Serializer\Annotation\SerializedName('lastUpdatedAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    public \DateTime $lastUpdatedAt;
    
	public function __construct()
	{
		$this->accountId = "";
		$this->balance = 0;
		$this->createdAt = new \DateTime();
		$this->currency = "";
		$this->lastUpdatedAt = new \DateTime();
	}
}
