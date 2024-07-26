<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class OrchestrationCreditWalletRequest
{
    #[\JMS\Serializer\Annotation\SerializedName('amount')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\Monetary')]
    public Monetary $amount;

    /**
     * The balance to credit
     *
     * @var ?string $balance
     */
    #[\JMS\Serializer\Annotation\SerializedName('balance')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $balance = null;

    /**
     * Metadata associated with the wallet.
     *
     * @var array<string, string> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;

    #[\JMS\Serializer\Annotation\SerializedName('reference')]
    #[\JMS\Serializer\Annotation\Type('string')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $reference = null;

    /**
     * $sources
     *
     * @var array<mixed> $sources
     */
    #[\JMS\Serializer\Annotation\SerializedName('sources')]
    #[\JMS\Serializer\Annotation\Type('array<mixed>')]
    public array $sources;

    #[\JMS\Serializer\Annotation\SerializedName('timestamp')]
    #[\JMS\Serializer\Annotation\Type("DateTime<'Y-m-d\TH:i:s.up'>")]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?\DateTime $timestamp = null;

    public function __construct()
    {
        $this->amount = new \formance\stack\Models\Shared\Monetary();
        $this->balance = null;
        $this->metadata = [];
        $this->reference = null;
        $this->sources = [];
        $this->timestamp = null;
    }
}