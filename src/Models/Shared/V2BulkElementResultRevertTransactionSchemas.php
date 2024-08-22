<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2BulkElementResultRevertTransactionSchemas
{
    /**
     *
     * @var V2Transaction $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\V2Transaction')]
    public V2Transaction $data;

    /**
     *
     * @var string $responseType
     */
    #[\JMS\Serializer\Annotation\SerializedName('responseType')]
    public string $responseType;

    /**
     * @param  ?V2Transaction  $data
     * @param  ?string  $responseType
     */
    public function __construct(?V2Transaction $data = null, ?string $responseType = null)
    {
        $this->data = $data;
        $this->responseType = $responseType;
    }
}