<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Payment
{
	#[\JMS\Serializer\Annotation\SerializedName('accountID')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $accountID;
    
    /**
     * $adjustments
     * 
     * @var array<\formance\stack\Models\Shared\PaymentAdjustment> $adjustments
     */
	#[\JMS\Serializer\Annotation\SerializedName('adjustments')]
    #[\JMS\Serializer\Annotation\Type('array<formance\stack\Models\Shared\PaymentAdjustment>')]
    public array $adjustments;
    
	#[\JMS\Serializer\Annotation\SerializedName('asset')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $asset;
    
	#[\JMS\Serializer\Annotation\SerializedName('createdAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    public \DateTime $createdAt;
    
	#[\JMS\Serializer\Annotation\SerializedName('id')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $id;
    
	#[\JMS\Serializer\Annotation\SerializedName('initialAmount')]
    #[\JMS\Serializer\Annotation\Type('int')]
    public int $initialAmount;
    
	#[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\PaymentMetadata')]
    public PaymentMetadata $metadata;
    
	#[\JMS\Serializer\Annotation\SerializedName('provider')]
    #[\JMS\Serializer\Annotation\Type('enum<formance\stack\Models\Shared\Connector>')]
    public Connector $provider;
    
    /**
     * $raw
     * 
     * @var array<string, mixed> $raw
     */
	#[\JMS\Serializer\Annotation\SerializedName('raw')]
    #[\JMS\Serializer\Annotation\Type('array<string, mixed>')]
    public array $raw;
    
	#[\JMS\Serializer\Annotation\SerializedName('reference')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $reference;
    
	#[\JMS\Serializer\Annotation\SerializedName('scheme')]
    #[\JMS\Serializer\Annotation\Type('enum<formance\stack\Models\Shared\PaymentScheme>')]
    public PaymentScheme $scheme;
    
	#[\JMS\Serializer\Annotation\SerializedName('status')]
    #[\JMS\Serializer\Annotation\Type('enum<formance\stack\Models\Shared\PaymentStatus>')]
    public PaymentStatus $status;
    
	#[\JMS\Serializer\Annotation\SerializedName('type')]
    #[\JMS\Serializer\Annotation\Type('enum<formance\stack\Models\Shared\PaymentType>')]
    public PaymentType $type;
    
	public function __construct()
	{
		$this->accountID = "";
		$this->adjustments = [];
		$this->asset = "";
		$this->createdAt = new \DateTime();
		$this->id = "";
		$this->initialAmount = 0;
		$this->metadata = new \formance\stack\Models\Shared\PaymentMetadata();
		$this->provider = \formance\stack\Models\Shared\Connector::STRIPE;
		$this->raw = [];
		$this->reference = "";
		$this->scheme = \formance\stack\Models\Shared\PaymentScheme::VISA;
		$this->status = \formance\stack\Models\Shared\PaymentStatus::PENDING;
		$this->type = \formance\stack\Models\Shared\PaymentType::PAY_IN;
	}
}
