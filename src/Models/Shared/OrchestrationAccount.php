<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class OrchestrationAccount
{
    /**
     *
     * @var string $address
     */
    #[\JMS\Serializer\Annotation\SerializedName('address')]
    public string $address;

    /**
     * $effectiveVolumes
     *
     * @var ?array<string, Volume> $effectiveVolumes
     */
    #[\JMS\Serializer\Annotation\SerializedName('effectiveVolumes')]
    #[\JMS\Serializer\Annotation\Type('array<string, \formance\stack\Models\Shared\Volume>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $effectiveVolumes = null;

    /**
     * $metadata
     *
     * @var array<string, string> $metadata
     */
    #[\JMS\Serializer\Annotation\SerializedName('metadata')]
    #[\JMS\Serializer\Annotation\Type('array<string, string>')]
    public array $metadata;

    /**
     * $volumes
     *
     * @var ?array<string, Volume> $volumes
     */
    #[\JMS\Serializer\Annotation\SerializedName('volumes')]
    #[\JMS\Serializer\Annotation\Type('array<string, \formance\stack\Models\Shared\Volume>')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?array $volumes = null;

    /**
     * @param  ?string  $address
     * @param  ?array<string, string>  $metadata
     * @param  ?array<string, Volume>  $effectiveVolumes
     * @param  ?array<string, Volume>  $volumes
     */
    public function __construct(?string $address = null, ?array $metadata = null, ?array $effectiveVolumes = null, ?array $volumes = null)
    {
        $this->address = $address;
        $this->metadata = $metadata;
        $this->effectiveVolumes = $effectiveVolumes;
        $this->volumes = $volumes;
    }
}