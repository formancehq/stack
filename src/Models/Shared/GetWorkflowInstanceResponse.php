<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class GetWorkflowInstanceResponse
{
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\WorkflowInstance')]
    public WorkflowInstance $data;

    public function __construct()
    {
        $this->data = new \formance\stack\Models\Shared\WorkflowInstance();
    }
}