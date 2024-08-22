<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class LedgerAccountSubject
{
    /**
     *
     * @var string $identifier
     */
    #[\JMS\Serializer\Annotation\SerializedName('identifier')]
    public string $identifier;

    /**
     *
     * @var string $type
     */
    #[\JMS\Serializer\Annotation\SerializedName('type')]
    public string $type;

    /**
     * @param  ?string  $identifier
     * @param  ?string  $type
     */
    public function __construct(?string $identifier = null, ?string $type = null)
    {
        $this->identifier = $identifier;
        $this->type = $type;
    }
}