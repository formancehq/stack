<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class Hold
{
    #[\JMS\Serializer\Annotation\SerializedName('description')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $description;

    #[\JMS\Serializer\Annotation\SerializedName('destination')]
    #[\JMS\Serializer\Annotation\Type('mixed')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public mixed $destination = null;

    /**
     * The unique ID of the hold.
     *
     * @var string $id
     */
    #[\JMS\Serializer\Annotation\SerializedName('id')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $id;

    /**
     * Metadata associated with the hold.
     *
     * @var array<string, string> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;

    /**
     * The ID of the wallet the hold is associated with.
     *
     * @var string $walletID
     */
    #[\JMS\Serializer\Annotation\SerializedName('walletID')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $walletID;

    public function __construct()
    {
        $this->description = '';
        $this->destination = null;
        $this->id = '';
        $this->metadata = [];
        $this->walletID = '';
    }
}