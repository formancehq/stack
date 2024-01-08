<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class PaymentsAccount
{
	#[\JMS\Serializer\Annotation\SerializedName('accountName')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $accountName;
    
	#[\JMS\Serializer\Annotation\SerializedName('connectorID')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $connectorID;
    
	#[\JMS\Serializer\Annotation\SerializedName('createdAt')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    public \DateTime $createdAt;
    
	#[\JMS\Serializer\Annotation\SerializedName('defaultAsset')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $defaultAsset;
    
    /**
     * 
     * @var string $defaultCurrency
     * @deprecated  field: This will be removed in a future release, please migrate away from it as soon as possible.
     */
	#[\JMS\Serializer\Annotation\SerializedName('defaultCurrency')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $defaultCurrency;
    
	#[\JMS\Serializer\Annotation\SerializedName('id')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $id;
    
    /**
     * $metadata
     * 
     * @var array<string, string> $metadata
     */
	#[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;
    
    /**
     * $pools
     * 
     * @var ?array<string> $pools
     */
	#[\JMS\Serializer\Annotation\SerializedName('pools')]
    #[\JMS\Serializer\Annotation\Type('array<string>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $pools = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('raw')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\PaymentsAccountRaw')]
    public PaymentsAccountRaw $raw;
    
	#[\JMS\Serializer\Annotation\SerializedName('reference')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $reference;
    
	#[\JMS\Serializer\Annotation\SerializedName('type')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $type;
    
	public function __construct()
	{
		$this->accountName = "";
		$this->connectorID = "";
		$this->createdAt = new \DateTime();
		$this->defaultAsset = "";
		$this->defaultCurrency = "";
		$this->id = "";
		$this->metadata = [];
		$this->pools = null;
		$this->raw = new \formance\stack\Models\Shared\PaymentsAccountRaw();
		$this->reference = "";
		$this->type = "";
	}
}
