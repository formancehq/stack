<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2StageStatus
{
	#[\JMS\Serializer\Annotation\SerializedName('error')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $error = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('instanceID')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $instanceID;
    
	#[\JMS\Serializer\Annotation\SerializedName('stage')]
    #[\JMS\Serializer\Annotation\Type('float')]
    public float $stage;
    
	#[\JMS\Serializer\Annotation\SerializedName('startedAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    public \DateTime $startedAt;
    
	#[\JMS\Serializer\Annotation\SerializedName('terminatedAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?\DateTime $terminatedAt = null;
    
	public function __construct()
	{
		$this->error = null;
		$this->instanceID = "";
		$this->stage = 0;
		$this->startedAt = new \DateTime();
		$this->terminatedAt = null;
	}
}
