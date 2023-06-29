<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class CreateScopeRequest
{
	#[\JMS\Serializer\Annotation\SerializedName('label')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $label;
    
    /**
     * $metadata
     * 
     * @var ?array<string, mixed> $metadata
     */
	#[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, mixed>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $metadata = null;
    
	public function __construct()
	{
		$this->label = "";
		$this->metadata = null;
	}
}
