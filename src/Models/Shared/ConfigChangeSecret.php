<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class ConfigChangeSecret
{
    #[\JMS\Serializer\Annotation\SerializedName('secret')]
    #[\JMS\Serializer\Annotation\Type('string')]
    public string $secret;

    public function __construct()
    {
        $this->secret = '';
    }
}