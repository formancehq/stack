<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class V2GetWorkflowResponse
{
    /**
     *
     * @var V2Workflow $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\V2Workflow')]
    public V2Workflow $data;

    /**
     * @param  ?V2Workflow  $data
     */
    public function __construct(?V2Workflow $data = null)
    {
        $this->data = $data;
    }
}